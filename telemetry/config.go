package telemetry

import (
	"time"
)

const (
	GCP        = "gcp"
	PROMETHEUS = "prometheus"
	STDOUT     = "stdout"
)

// MetricsConfig struct
type MetricsConfig struct {
	Enabled  bool   `yaml:",omitempty" default:"false"`
	Backend  string `yaml:",omitempty" validate:"oneof=gcp prometheus stdout" default:"stdout"`
	Endpoint string `yaml:"endpoint,omitempty"`
	// SamplingFraction >= 1 will always sample. SamplingFraction < 0 are treated as zero.
	SamplingFraction float64       `yaml:"sampling_fraction,omitempty" default:"1.0"`
	CollectPeriod    time.Duration `yaml:"collect_period,omitempty" default:"10s"`
    HistogramDistribution []float64 `yaml:"histogram_distribution,omitempty" validate:"omitempty,gt=0,dive,min=0,max=10"`
}

// TracingConfig struct
type TracingConfig struct {
	Enabled  bool   `yaml:",omitempty" default:"false"`
	Backend  string `yaml:",omitempty" validate:"oneof=gcp stdout" default:"stdout"`
	Endpoint string `yaml:"endpoint,omitempty"`
	// SamplingFraction >= 1 will always sample. SamplingFraction < 0 are treated as zero.
	SamplingFraction float64 `yaml:"sampling_fraction,omitempty" default:"1.0"`
}
