# Telemetry

**OpenTelemetry** initialization helpers



## Usage

Optional environment variables 

```bash
export OTEL_RESOURCE_ATTRIBUTES=key=value,rk5=7
export OTEL_SERVICE_NAME=play-service

# enable metrics and tracking
export CONFY_FEATURES_METRICS_ENABLED=true
export CONFY_FEATURES_TRACING_ENABLED=true
# enable metrics target: `prometheus` and tracing target: `stdout`
export CONFY_FEATURES_METRICS_TARGET=prometheus
export CONFY_FEATURES_TRACING_TARGET=stdout
# when using with target: `gcp`
export GOOGLE_CLOUD_PROJECT=xyz
export GOOGLE_APPLICATION_CREDENTIALS=../../../Apps/micro-starter-kit.json
```

```go
import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/instrumentation/grpctrace"

    "github.com/xmlking/toolkit/telemetry/metrics"
    "github.com/xmlking/toolkit/telemetry/tracing"
)

func main() {
    cfg := config.GetConfig()
    
    if cfg.Features.Tracing.Enabled {
    closeFn := tracing.InitTracing(ctx, cfg.Features.Tracing)
    defer closeFn()
    }

    if cfg.Features.Metrics.Enabled {
    closeFn := metrics.InitMetrics(ctx, cfg.Features.Metrics)
    defer closeFn()
    }

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctrace.UnaryServerInterceptor(global.Tracer(""))),
		grpc.StreamInterceptor(grpctrace.StreamServerInterceptor(global.Tracer(""))),
	)

	api.RegisterHelloServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

### Examples

#### Integrated GCP observability

- https://cloud.google.com/trace/docs/setup/go-ot#gke
- https://medium.com/google-cloud/integrating-tracing-and-logging-with-opentelemetry-and-stackdriver-a5396fbc3e78

#### gRPC tracking

- https://github.com/open-telemetry/opentelemetry-go/blob/master/example/grpc/server/main.go
- 
- https://github.com/open-telemetry/opentelemetry-go/blob/master/exporters/stdout/example_test.go

#### Metrics:

- https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
- https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/blob/master/example/metric/example.go
- https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/tree/main/example/metric
- https://opentelemetry.uptrace.dev/guide/metrics.html#synchronous-instruments
