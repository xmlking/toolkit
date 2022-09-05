package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server is a grpc Server
type Server interface {
	// Server returns grpc server previously created with NewServer()
	Server() *grpc.Server
	// Client creates and returns a new grpc client connection
	// Users should call ClientConn.Close() to terminate all the pending operations.
	// also keep track of all clients so that, it disconnect client connections gracefully on kill signal.
	// throw error if connection fail or if there is already a client created with same name
	Client(target string, opts ...ClientOption) (*grpc.ClientConn, error)
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

// NewClient creates and returns a new grpc client connection
// Users should call ClientConn.Close() to terminate all the pending operations.
// also keep track of all clients so that, it disconnects client connections gracefully on kill signal.
// throw error if connection fail or if there is already a client created with same name
func NewClient(target string, opts ...ClientOption) (*grpc.ClientConn, error) {
	return DefaultServer.Client(target, opts...)
}

// SetServingStatus updates service health status. concurrency safe
func SetServingStatus(service string, servingStatus grpc_health_v1.HealthCheckResponse_ServingStatus) {
	DefaultServer.SetServingStatus(service, servingStatus)
}

// Start all daemon services and wait for kill signal, then gracefully shutdown
func Start() error {
	return DefaultServer.Start()
}
