package service

import (
	"context"

	"google.golang.org/grpc"
)

type Option func(*Options)

type Options struct {
	Name    string
	Version string

	GrpcEndpoint string
	GrpcOptions  []grpc.ServerOption
	DialOptions  []grpc.DialOption

	Context context.Context

	// Before and After funcs
	//BeforeStart []func() error
	//BeforeStop  []func() error
	//AfterStart  []func() error
	//AfterStop   []func() error

}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func WithGrpcOptions(opts ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.GrpcOptions = opts
	}
}

func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(o *Options) {
		o.DialOptions = opts
	}
}

// WithGrpcEndpoint specifies the net.Listener endpoint to use instead of the default
func WithGrpcEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.GrpcEndpoint = endpoint
	}
}
