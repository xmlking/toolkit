package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type Remote struct {
	Endpoint      string
	ServiceConfig string
	Authority     string
}

// Server is a grpc Server
type Server interface {
	NewClient(remote Remote) (*grpc.ClientConn, error)
	Start() error
	Stop()
}

// NewServer creates and returns a new Broker based on the packages within.
func NewServer(ctx context.Context, opts ...Option) Server {
	return newServer(ctx, opts...)
}
