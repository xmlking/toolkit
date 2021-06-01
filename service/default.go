package service

import (
	"context"
	"net"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/xmlking/toolkit/util/endpoint"
)

const (
	DefaultName    = "mkit.service.default"
	DefaultVersion = "latest"
	DefaultAddress = ":0"
	// DefaultShutdownTimeout defines the default timeout given to the service when calling Shutdown.
	DefaultShutdownTimeout = time.Minute * 1
)

type service struct {
	opts       Options
	grpcServer *grpc.Server
}

func newService(opts ...Option) Service {
	// Default Options
	options := Options{
		Name:    DefaultName,
		Version: DefaultVersion,
	}
	s := service{opts: options}
	s.ApplyOptions(opts...)

	s.grpcServer = grpc.NewServer(s.opts.GrpcOptions...)

	return &s
}

func (s *service) ApplyOptions(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Server() *grpc.Server {
	return s.grpcServer
}

func (s *service) Client(remote Remote) (clientConn *grpc.ClientConn, err error) {
	var dialOptions []grpc.DialOption

	// TODO: set TLS also

	if remote.ServiceConfig != "" {
		dialOptions = append(s.opts.DialOptions, grpc.WithDefaultServiceConfig(remote.ServiceConfig))
	}
	if remote.Authority != "" {
		dialOptions = append(s.opts.DialOptions, grpc.WithAuthority(remote.Authority))
	}

	clientConn, err = grpc.Dial(remote.Endpoint, dialOptions...)
	if err != nil {
		log.Error().Err(err).Msgf("Failed connect to: %s", remote.Endpoint)
	}
	return
}

func (s *service) Shutdown() error {
	return nil
}

func (s *service) Start() (err error) {
	//println(config.GetBuildInfo())

	// eg, egCtx := errgroup.WithContext(s.opts.Context)
	ctx, cancel := context.WithCancel(s.opts.Context)
	defer cancel()

	errCh := make(chan error, 1)

	// Start GrpcServer
	// Add HealthChecks
	hsrv := health.NewServer()
	for name := range s.grpcServer.GetServiceInfo() {
		hsrv.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(s.grpcServer, hsrv)
	// TODO: User our own custom health implementation, instead of using built-in health server
	// https://github.com/GoogleCloudPlatform/grpc-gke-nlb-tutorial/blob/master/echo-grpc/health/health.go

	var listener net.Listener
	if s.opts.GrpcEndpoint == "" {
		listener, err = net.Listen("tcp", DefaultAddress)
		if err != nil {
			return errors.Wrap(err, "Failed to create listener")
		}
	} else {
		listener, err = endpoint.GetListener(s.opts.GrpcEndpoint)
		if err != nil {
			return errors.Wrap(err, "Failed to create listener")
		}
	}
	// log.Info().Msgf("Server (%s) starting at: %s, secure: %t", s.opts.Name, listener.Addr(), s.cfg.Features.Tls.Enabled)
	log.Info().Msgf("Server (%s) starting at: %s", s.opts.Name, listener.Addr())
	go func() {
		reflection.Register(s.grpcServer)
		errCh <- s.grpcServer.Serve(listener)
	}()

	// This will block until either a signal arrives or one of the grouped functions
	// returns an error.
	// <-egCtx.Done()

	// Stop either if the receiver stops (sending to errCh) or if stopCh is closed.
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		break
	}

	// do any more resources closing here
	s.grpcServer.GracefulStop()
	cancel()

	select {
	case <-errCh:
		s.grpcServer.Stop()
		log.Info().Msg("Gracefully shutdown")
	case <-time.After(DefaultShutdownTimeout):
		log.Error().Msg("Failed to shutdown within grace period")
		return errors.New("Failed to shutdown within grace period")
	}

	return nil
}
