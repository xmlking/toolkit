package grpc

import (
	"context"
	"net"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/util/endpoint"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	DefaultAddress = ":0"
)

type grpcServer struct {
	options Options
	server  *grpc.Server
}

func newServer(ctx context.Context, opts ...Option) Server {
	// Default Options
	options := Options{
		Context: ctx,
	}

	for _, o := range opts {
		o(&options)
	}

	//interceptors := []grpc.UnaryServerInterceptor{
	//    // interceptor.NewAuthTokenPropagator(),
	//}
	//
	//opts := []grpc.ServerOption{
	//    middleware.WithUnaryServerChain(interceptors...),
	//}

	server := grpc.NewServer(options.GrpcOptions...)
	reflection.Register(server)

	return &grpcServer{
		options: options,
		server:  server,
	}
}

func (s *grpcServer) NewClient(remote Remote) (clientConn *grpc.ClientConn, err error) {
	var dialOptions []grpc.DialOption

	// TODO: set TLS also

	if remote.ServiceConfig != "" {
		dialOptions = append(s.options.DialOptions, grpc.WithDefaultServiceConfig(remote.ServiceConfig))
	}
	if remote.Authority != "" {
		dialOptions = append(s.options.DialOptions, grpc.WithAuthority(remote.Authority))
	}

	clientConn, err = grpc.Dial(remote.Endpoint, dialOptions...)
	if err != nil {
		log.Error().Err(err).Msgf("Failed connect to: %s", remote.Endpoint)
	}
	return
}

func (s *grpcServer) Start() (err error) {
	ctx := s.options.Context
	g, _ := errgroup.WithContext(s.options.Context)

	// Add HealthChecks
	hsrv := health.NewServer()
	for name := range s.server.GetServiceInfo() {
		hsrv.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(s.server, hsrv)
	// TODO: User our own custom health implementation, instead of using built-in health server
	// https://github.com/GoogleCloudPlatform/grpc-gke-nlb-tutorial/blob/master/echo-grpc/health/health.go

	var listener net.Listener
	if s.options.GrpcEndpoint == "" {
		listener, err = net.Listen("tcp", DefaultAddress)
		if err != nil {
			return errors.Wrap(err, "Failed to create listener")
		}
	} else {
		listener, err = endpoint.GetListener(s.options.GrpcEndpoint)
		if err != nil {
			return errors.Wrap(err, "Failed to create listener")
		}
	}

	// log.Info().Msgf("Server starting at: %s, secure: %t", listener.Addr(), s.cfg.Features.Tls.Enabled)
	log.Info().Msgf("Server starting at: %s", listener.Addr())
	g.Go(func() error {
		return s.server.Serve(listener)
	})

	g.Go(func() (err error) {
		// listen for the interrupt signal
		<-ctx.Done()

		// log situation
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Debug().Str("component", "grpc").Msg("Context timeout exceeded")
		case context.Canceled:
			log.Debug().Str("component", "grpc").Msg("Context cancelled by interrupt signal")
		}

		// wait for all subs to stop
		log.Info().Str("component", "grpc").Msgf("Closing grpc server...")
		s.server.GracefulStop()

		return
	})

	// Wait for all tasks to be finished or return if error occur at any task.
	return g.Wait()
}

func (s *grpcServer) Stop() {
	s.server.Stop()
}
