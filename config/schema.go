package config

import (
	"time"
)

type Service struct {
	Endpoint      string            `yaml:"endpoint" required:"true"`
	Version       string            `yaml:",omitempty" default:"v0.1.0"`
	Metadata      map[string]string `yaml:"metadata,omitempty"`
	ServiceConfig string            `yaml:"service_config,omitempty"`
	Authority     string            `yaml:"authority,omitempty"`
}

// Example service
type Services struct {
	Account *Service `yaml:"account,omitempty"`
}

type Features struct {
	Metrics   *Features_Metrics   `yaml:"metrics,omitempty"`
	Tracing   *Features_Tracing   `yaml:"tracing,omitempty"`
	TLS       *Features_TLS       `yaml:"tls,omitempty"`
	Validator *Features_Validator `yaml:"validator,omitempty"`
	Rpclog    *Features_Rpclog    `yaml:"rpclog,omitempty"`
	Translog  *Features_Translog  `yaml:"translog,omitempty"`
}

type Features_Metrics struct {
	Enabled       bool   `yaml:",omitempty" default:"false"`
	Address       string `yaml:"address,omitempty"`
	FlushInterval uint64 `yaml:"flush_interval,omitempty" default:"10000000"`
}

type Features_Tracing struct {
	Enabled       bool   `yaml:",omitempty" default:"false"`
	Address       string `yaml:"address,omitempty"`
	FlushInterval uint64 `yaml:"flush_interval,omitempty" default:"10000000"`
}

type Features_TLS struct {
	Enabled    bool   `yaml:",omitempty" default:"false"`
	CertFile   string `yaml:"cert_file" valid:"type(string),required"`
	KeyFile    string `yaml:"key_file" valid:"type(string),required"`
	CaFile     string `yaml:"ca_file" valid:"type(string),required"`
	Password   string `yaml:"password,omitempty"`
	ServerName string `yaml:"server_name,omitempty" default:"'*'"`
}

type Features_Validator struct {
	Enabled bool `yaml:",omitempty" default:"false"`
}

type Features_Rpclog struct {
	Enabled bool `yaml:",omitempty" default:"false"`
}

type Features_Translog struct {
	Enabled bool   `yaml:",omitempty" default:"false"`
	Topic   string `yaml:"topic,omitempty"`
}

type PublishSettings struct {
	DelayThreshold    time.Duration `yaml:"delay_threshold,omitempty"`
	CountThreshold    int           `yaml:"count_threshold,omitempty"`
	ByteThreshold     int           `yaml:"byte_threshold,omitempty"`
	NumGoroutines     int           `yaml:"num_goroutines,omitempty"`
	Timeout           time.Duration `yaml:",omitempty"`
	BufferedByteLimit int           `yaml:"buffered_byte_limit,omitempty"`
}

type ReceiveSettings struct {
	MaxExtension           time.Duration `yaml:"max_extension,omitempty"`
	MaxExtensionPeriod     time.Duration `yaml:"max_extension_period,omitempty"`
	MaxOutstandingMessages int           `yaml:"max_outstanding_messages,omitempty"`
	MaxOutstandingBytes    int           `yaml:"max_outstanding_bytes,omitempty"`
	NumGoroutines          int           `yaml:"num_goroutines,omitempty"`
	Synchronous            bool          `yaml:",omitempty" default:"false"`
}

type Configuration struct {
	ProjectID string    `yaml:"project_id,omitempty" env:"GOOGLE_CLOUD_PROJECT" default:"my-project-id"`
	Features  *Features `yaml:"features,omitempty"`
	Services  *Services `yaml:"services,omitempty"`
}
