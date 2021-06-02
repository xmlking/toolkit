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
	gSrv    *grpc.Server
	hSrv    *health.Server
	clients map[string]*grpc.ClientConn
}

func (s *grpcServer) SetServingStatus(service string, servingStatus grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.hSrv.SetServingStatus(service, servingStatus)
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

	// Add HealthChecks only after all user services are registered
	s.hSrv = health.NewServer()
	for name := range s.gSrv.GetServiceInfo() {
		s.hSrv.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(s.gSrv, s.hSrv)

	// registers the server reflection service on the given gRPC server.
	reflection.Register(s.gSrv)

	g.Go(func() error {
		return s.gSrv.Serve(s.options.Listener)
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

		// Gracefully stop healthServer
		s.hSrv.Shutdown()

		// Gracefully stop clients
		for name, client := range s.clients {
			log.Info().Str("component", "grpc").Msgf("Closing grpc client(%s) for %s", name, client.Target())
			err = client.Close()
		}
		log.Info().Str("component", "grpc").Msg("Stopping grpc server...")
		// Gracefully stop server
		s.gSrv.GracefulStop()

		return
	})

	// Wait for all tasks to be finished or return if error occur at any task.
	return g.Wait()
}

func (s *grpcServer) Stop() {
	s.gSrv.Stop()
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
		gSrv:    server,
		clients: make(map[string]*grpc.ClientConn),
	}
}
