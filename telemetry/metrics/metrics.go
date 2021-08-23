package metrics

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"

	cloudmetric "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	goprom "github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	once sync.Once
	cont *controller.Controller
)

// InitMetrics Initialize Metrics exporter
// InitMetrics expected GOOGLE_CLOUD_PROJECT & GOOGLE_APPLICATION_CREDENTIALS Environment Variable set
// Usage: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func InitMetrics(ctx context.Context, cfg *telemetry.MetricsConfig) func() {
	once.Do(func() {
		log.Debug().Interface("MetricConfig", cfg).Msg("Initializing Metrics")

		resources, err := resource.New(ctx,
			// Builtin detectors provide default values and support
			// OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
			resource.WithProcess(), // This option configures a set of Detectors that discover process information
			// resource.WithAttributes(attribute.String("foo", "bar")), // Or specify resource attributes directly
		)
		if err != nil {
			log.Fatal().Stack().Err(err).Msg("failed to initialize resources for metrics exporter")
		}

		switch cfg.Backend {
		case telemetry.GCP:
			projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
			opts := []cloudmetric.Option{cloudmetric.WithProjectID(projectID)}
			popts := []controller.Option{
				controller.WithCollectPeriod(cfg.CollectPeriod),
				controller.WithResource(resources),
			}
			cont, err = cloudmetric.InstallNewPipeline(opts, popts...)
			if err != nil {
				log.Fatal().Stack().Err(err).Msg("failed to initialize metrics exporter")
			}

		case telemetry.PROMETHEUS:
			var pConfig = prometheus.Config{Registry: goprom.DefaultRegisterer.(*goprom.Registry)}

			if cfg.HistogramDistribution != nil {
				pConfig.DefaultHistogramBoundaries = cfg.HistogramDistribution
			}

			pController := controller.New(
				processor.New(
					simple.NewWithHistogramDistribution(
						histogram.WithExplicitBoundaries(pConfig.DefaultHistogramBoundaries),
					),
					export.CumulativeExportKindSelector(),
					processor.WithMemory(true),
				),
				controller.WithCollectPeriod(cfg.CollectPeriod),
				controller.WithResource(resources),
			)

			exporter, err := prometheus.New(pConfig, pController)
			cont = exporter.Controller()

			if err != nil {
				log.Fatal().Stack().Err(err).Msg("failed to initialize prometheus exporter")
			}

			// start prometheus exporter
			http.HandleFunc("/metrics", exporter.ServeHTTP)
			pSrv := &http.Server{
				// use appCtx which get canceled on signal
				BaseContext: func(_ net.Listener) context.Context { return ctx },
			}
			listener, err := net.Listen("tcp", cfg.Endpoint)
			if err != nil {
				log.Fatal().Stack().Err(err).Msg("error creating listener for prometheus exporter")
			}
			go func() {
				if err := pSrv.Serve(listener); err != http.ErrServerClosed {
					log.Fatal().Stack().Err(err).Msg("Prometheus exporter error:")
				}
			}()

			log.Info().Msgf("Prometheus exporter running at: %s\n", listener.Addr())

		case telemetry.STDOUT:
			opts := []stdoutmetric.Option{
				stdoutmetric.WithPrettyPrint(),
			}
			var metricExporter *stdoutmetric.Exporter
			metricExporter, err = stdoutmetric.New(opts...)
			if err != nil {
				log.Fatal().Stack().Err(err).Msg("failed to initialize metrics exporter")
			}
			cont = controller.New(
				processor.New(
					simple.NewWithExactDistribution(),
					metricExporter,
				),
				controller.WithExporter(metricExporter),
				controller.WithCollectPeriod(cfg.CollectPeriod),
				controller.WithResource(resources),
			)
			err = cont.Start(ctx)
			if err != nil {
				log.Fatal().Stack().Err(err).Msg("failed to initialize metrics controller")
			}

		default:
			log.Fatal().Msgf("unsupported tracing Backend: '%s'", cfg.Backend)
		}

		// Registers metrics Provider globally.
		global.SetMeterProvider(cont.MeterProvider())
		propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		otel.SetTextMapPropagator(propagator)
	})

	return func() {
		cont.Stop(ctx)
	}
}
