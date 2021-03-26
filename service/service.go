package service

import (
	"google.golang.org/grpc"

	"github.com/xmlking/toolkit/broker/pubsub"
)

type Remote struct {
	Endpoint      string
	ServiceConfig string
	Authority     string
}

type Service interface {
	Options() Options
	Server() *grpc.Server
	Client(remote Remote) (*grpc.ClientConn, error)
	Broker() broker.Broker
	ApplyOptions(opts ...Option) // TODO: no use, make private ?
	Start() error
	Shutdown() error
}

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return newService(opts...)
}
