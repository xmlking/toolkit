package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type ServerOption func(*ServerOptions)

type ServerOptions struct {
	Name          string
	Listener      net.Listener
	ServerOptions []grpc.ServerOption
	Context       context.Context
}

// ServerName Name of the service
func ServerName(n string) ServerOption {
	return func(o *ServerOptions) {
		o.Name = n
	}
}

// WithServerOptions ServerOptions for server
func WithServerOptions(opts ...grpc.ServerOption) ServerOption {
	return func(o *ServerOptions) {
		o.ServerOptions = opts
	}
}

// WithListener specifies the net.Listener to use instead of the default
func WithListener(listener net.Listener) ServerOption {
	return func(o *ServerOptions) {
		o.Listener = listener
	}
}

// Context  appContext to trigger terminate signal
func Context(ctx context.Context) ServerOption {
	return func(o *ServerOptions) {
		o.Context = ctx
	}
}
