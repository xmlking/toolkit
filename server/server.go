package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server is a grpc Server
type Server interface {
	// NewClient creates and returns a new grpc Client
	NewClient(target string, opts ...ClientOption) (*grpc.ClientConn, error)
	// SetServingStatus updates service health status. concurrency safe
	// empty service string represents the health of the whole system
	SetServingStatus(service string, servingStatus grpc_health_v1.HealthCheckResponse_ServingStatus)
	Start() error
	// Stop will force stop server
	Stop()
}

var DefaultServer Server

// NewServer creates and returns a new grpc Server
func NewServer(ctx context.Context, opts ...ServerOption) Server {
	DefaultServer = newServer(ctx, opts...)
	return DefaultServer
}

// NewClient creates and returns a new grpc Client
// Users should call ClientConn.Close() to terminate all the pending operations.
func NewClient(target string, opts ...ClientOption) (*grpc.ClientConn, error) {
	return DefaultServer.NewClient(target, opts...)
}

// SetServingStatus updates service health status. concurrency safe
func SetServingStatus(service string, servingStatus grpc_health_v1.HealthCheckResponse_ServingStatus) {
	DefaultServer.SetServingStatus(service, servingStatus)
}

func Start() error {
	return DefaultServer.Start()
}
