module github.com/xmlking/toolkit/telemetry

go 1.16

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v1.0.0-RC1
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.0.0-RC2
	github.com/rs/zerolog v1.23.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/prometheus v0.22.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.22.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC2
	go.opentelemetry.io/otel/metric v0.22.0
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/sdk/export/metric v0.22.0
	go.opentelemetry.io/otel/sdk/metric v0.22.0
)
