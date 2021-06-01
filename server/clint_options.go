package server

import (
	"context"

	"google.golang.org/grpc"
)

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	Name        string
	DialOptions []grpc.DialOption
	Context     context.Context
}

// ClientName Name of the client
func ClientName(n string) ClientOption {
	return func(o *ClientOptions) {
		o.Name = n
	}
}

// WithDialOptions DialOptions for client
func WithDialOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *ClientOptions) {
		o.DialOptions = opts
	}
}
