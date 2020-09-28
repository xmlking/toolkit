package service

import (
	"google.golang.org/grpc"

	"github.com/xmlking/toolkit"
)

type Remote struct {
	Endpoint      string
	ServiceConfig string
	Authority     string
}

type Service interface {
	Server() *grpc.Server
	Client(remote Remote) (*grpc.ClientConn, error)
	Options() Options
	ApplyOptions(opts ...Option) // TODO: no use, make private ?
	AddSubscriber(fn interface{})
	// Stop the service
	Shutdown() error
	// Run the service
	Start() error
	//Config Interface
}

// TODO rename Service --> GrpcService
type GrpcService interface {
	toolkit.Service
	GrpcServer() *grpc.Server
	GrpcClient(remote Remote) (*grpc.ClientConn, error)
	Options() Options
	ApplyOptions(opts ...Option) // TODO: no use, make private ?
}

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return newService(opts...)
}
