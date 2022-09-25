package types

import (
	"time"
)

/** Feature Types **/

// ServerTLS struct
type ServerTLS struct {
	Enabled    bool   `yaml:",omitempty" default:"false"`
	CertFile   string `yaml:"cert_file,omitempty" valid:"type(string),required_if=Enabled true"` // validate:"file,required_if=Enabled true"
	KeyFile    string `yaml:"key_file,omitempty" valid:"type(string),required_if=Enabled true"`  // validate:"file,required_if=Enabled true"
	CaFile     string `yaml:"ca_file,omitempty" valid:"type(string),required_if=Enabled true"`   // validate:"file,required_if=Enabled true"
	VerifyPeer bool   `yaml:"verify_peer,omitempty" default:"false"`
}

// ClientTLS struct
type ClientTLS struct {
	Enabled    bool   `yaml:",omitempty" default:"false"`
	CertFile   string `yaml:"cert_file,omitempty" valid:"type(string),required_if=Enabled true"` // validate:"file,required_if=Enabled true"
	KeyFile    string `yaml:"key_file,omitempty" valid:"type(string),required_if=Enabled true"`  // validate:"file,required_if=Enabled true"
	CaFile     string `yaml:"ca_file,omitempty" valid:"omitempty,type(string)"`                  // validate:"file,required_if=Enabled true"
	ServerName string `yaml:"server_name,omitempty" default:"'*'"`                               // Server Name Indication (SNI) aka Authority
}

// Validator struct
type Validator struct {
	// Feature
	Enabled bool `yaml:",omitempty" default:"false"`
}

// Rpclog struct
type Rpclog struct {
	Enabled bool `yaml:",omitempty" default:"false"`
}

// Translog struct
type Translog struct {
	Enabled bool   `yaml:",omitempty" default:"false"`
	Topic   string `yaml:",omitempty"`
}

// HealthCheck struct
type HealthCheck struct {
	Enabled  bool          `yaml:",omitempty" default:"false"`
	Interval time.Duration `yaml:",omitempty" validate:"omitempty,gte=1s" default:"60s"`
}

// PublishSettings : GCP PubSub Settings
type PublishSettings struct {
	Enabled           bool          `yaml:",omitempty" default:"false"`
	DelayThreshold    time.Duration `yaml:"delay_threshold,omitempty"`
	CountThreshold    int           `yaml:"count_threshold,omitempty"`
	ByteThreshold     int           `yaml:"byte_threshold,omitempty"`
	NumGoroutines     int           `yaml:"num_goroutines,omitempty"`
	Timeout           time.Duration `yaml:",omitempty"`
	BufferedByteLimit int           `yaml:"buffered_byte_limit,omitempty"`
}

// ReceiveSettings : GCP PubSub Settings
type ReceiveSettings struct {
	Enabled                bool          `yaml:",omitempty" default:"false"`
	MaxExtension           time.Duration `yaml:"max_extension,omitempty"`
	MaxExtensionPeriod     time.Duration `yaml:"max_extension_period,omitempty"`
	MaxOutstandingMessages int           `yaml:"max_outstanding_messages,omitempty"`
	MaxOutstandingBytes    int           `yaml:"max_outstanding_bytes,omitempty"`
	NumGoroutines          int           `yaml:"num_goroutines,omitempty"`
	Synchronous            bool          `yaml:",omitempty" default:"false"`
}

// Features : Example Features
type Features struct {
	//Metrics     *telemetry.MetricsConfig `yaml:"metrics,omitempty"`
	//Tracing     *telemetry.TracingConfig `yaml:"tracing,omitempty"`
	ServerTLS   *ServerTLS   `yaml:"server_tls,omitempty"`
	ClientTLS   *ClientTLS   `yaml:"client_tls,omitempty"`
	Validator   *Validator   `yaml:"validator,omitempty"`
	Rpclog      *Rpclog      `yaml:"rpclog,omitempty"`
	Translog    *Translog    `yaml:"translog,omitempty"`
	HealthCheck *HealthCheck `yaml:"health_check,omitempty"`
}

// Service config struct
type Service struct {
	Endpoint      string            `yaml:"endpoint" required:"true"`
	Version       string            `yaml:",omitempty" default:"v0.1.0"`
	Metadata      map[string]string `yaml:",omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ServiceConfig string            `yaml:"service_config,omitempty"`
	Authority     string            `yaml:",omitempty"`
}

// Services Example service
type Services struct {
	Account *Service `yaml:"account,omitempty"`
	Greeter *Service `yaml:"greeter,omitempty"`
	Emailer *Service `yaml:"emailer,omitempty"`
}

// Configuration : Example configuration
type Configuration struct {
	Host     string    `yaml:",omitempty" default:"0.0.0.0" validate:"ip"`
	Port     uint32    `yaml:",omitempty" default:"8080" validate:"numeric,gt=1024,lte=65535"`
	Features *Features `yaml:"features,omitempty"`
	Services *Services `yaml:"services,omitempty"`
	// ProjectID string    `yaml:"project_id,omitempty" env:"GOOGLE_CLOUD_PROJECT" default:"my-project-id"`
}
