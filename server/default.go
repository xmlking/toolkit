package server

import (
	"context"
	"net"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// TODO WithTracing() withMetrics() enable/disable

const (
	// DefaultName default server name
	DefaultName = "mkit.service.default"
	// DefaultAddress random port
	DefaultAddress = ":0"
	// DefaultMaxRecvMsgSize maximum message that client can receive
	// (16 MB).
	DefaultMaxRecvMsgSize = 1024 * 1024 * 16

	// DefaultMaxSendMsgSize maximum message that client can send
	// (16 MB).
	DefaultMaxSendMsgSize = 1024 * 1024 * 16
)

type grpcServer struct {
	options ServerOptions
	server  *grpc.Server
	clients map[string]*grpc.ClientConn
	status  chan grpc_health_v1.HealthCheckResponse_ServingStatus
}

func (s *grpcServer) SetServingStatus(status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.status <- status
}

func (s *grpcServer) NewClient(target string, opts ...ClientOption) (clientConn *grpc.ClientConn, err error) {
	// Default Options
	options := ClientOptions{
		Context: s.options.Context,
	}

	for _, o := range opts {
		o(&options)
	}

	// if client name not set, uses target as name
	if options.Name == "" {
		options.Name = target
	}

	if _, ok := s.clients[options.Name]; ok {
		return nil, errors.Newf("client with name: %s already created", options.Name)
	}

	if clientConn, err = grpc.DialContext(s.options.Context, target, options.DialOptions...); err != nil {
		log.Error().Err(err).Msgf("Failed connect to target: %s", target)
		return nil, err
	}

	s.clients[options.Name] = clientConn
	return
}

func (s *grpcServer) Start() (err error) {
	ctx := s.options.Context
	g, _ := errgroup.WithContext(s.options.Context)

	// Add HealthChecks after all user services are registered
	hsrv := health.NewServer()
	for name := range s.server.GetServiceInfo() {
		hsrv.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(s.server, hsrv)

	// asynchronously inspect dependencies and toggle serving status as needed
	go func() {
		for next := range s.status {
			log.Info().Str("component", "grpc").Msgf("Health status changed to: %s", next)
			// empty string represents the health of the whole system
			hsrv.SetServingStatus("", next)
		}
		log.Info().Str("component", "grpc").Msg("Stopped health status update watch job")
	}()

	// registers the server reflection service on the given gRPC server.
	reflection.Register(s.server)

	g.Go(func() error {
		return s.server.Serve(s.options.Listener)
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

		// Inform readiness prob to stop sending traffic
		SetServingStatus(grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		close(s.status)

		// Gracefully stop clients
		for name, client := range s.clients {
			log.Info().Str("component", "grpc").Msgf("Closing grpc client(%s) for %s", name, client.Target())
			err = client.Close()
		}
		log.Info().Str("component", "grpc").Msg("Stopping grpc server...")
		// Gracefully stop server
		s.server.GracefulStop()

		return
	})

	// Wait for all tasks to be finished or return if error occur at any task.
	return g.Wait()
}

func (s *grpcServer) Stop() {
	s.server.Stop()
}

func newServer(ctx context.Context, opts ...ServerOption) Server {
	// Default Options
	options := ServerOptions{
		Name:    DefaultName,
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

	if options.Listener == nil {
		var err error
		if options.Listener, err = net.Listen("tcp", DefaultAddress); err != nil {
			log.Fatal().Stack().Err(err).Str("component", "grpc").Msgf("Failed to create listener on DefaultAddress: %s", DefaultAddress)
		}
	}

	server := grpc.NewServer(options.ServerOptions...)

	return &grpcServer{
		options: options,
		server:  server,
		clients: make(map[string]*grpc.ClientConn),
		status:  make(chan grpc_health_v1.HealthCheckResponse_ServingStatus),
	}
}
