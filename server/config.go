package server

import (
	"go.opencensus.io/trace"
)

type (
	// ServerConfiguration is for gRPC server config.
	ServerConfiguration struct {
		Name          string
		Log           string
		Mode          string       `json:",default=pro,options=dev|test|rt|pre|pro"`
		MetricsURL    string       `json:",optional"`
		Telemetry     trace.Config `json:",optional"`
		ListenOn      string
		Auth          bool `json:",optional"`
		StrictControl bool `json:",optional"`
		// setting 0 means no timeout
		Timeout      int64 `json:",default=2000"`
		CPUThreshold int64 `json:",default=900,range=[0:1000]"`
	}

	// ClientConfiguration is for gRPC client config.
	ClientConfiguration struct {
		Endpoints []string `json:",optional"`
		Target    string   `json:",optional"`
		App       string   `json:",optional"`
		Token     string   `json:",optional"`
		NonBlock  bool     `json:",optional"`
		Timeout   int64    `json:",default=2000"`
	}
)
