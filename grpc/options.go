package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type Options struct {
	GrpcEndpoint string
	GrpcOptions  []grpc.ServerOption
	DialOptions  []grpc.DialOption
	Context      context.Context
}

type Option func(*Options)

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
