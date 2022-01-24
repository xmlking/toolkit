package xds

import "time"

const (
	STATIC     = "static"
	FILE       = "file"
	DNS        = "dns"
	KUBERNETES = "kubernetes"
)

// Configuration is for xds config
type Configuration struct {
	// MaxConcurrentStreams max concurrent streams for gRPC server
	MaxConcurrentStreams uint32        `yaml:"max_concurrent_streams,omitempty" validate:"omitempty,number" default:"1000000"`
	SourceType           string        `yaml:"source_type,omitempty" validate:"oneof=dns kubernetes file static" default:"static"`
	NodeID               string        `yaml:"node_id,omitempty" validate:"required,uuid"`
	RefreshInterval      time.Duration `yaml:"refresh_interval,omitempty" validate:"omitempty,min=0s,max=1h" default:"0s"`
	// Namespace to monitor when SourceType = kubernetes
	Namespace string `yaml:"namespace,omitempty" validate:"omitempty,alphanum" default:"xds"`
}
