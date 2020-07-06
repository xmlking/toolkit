# toolkit

## Usage

```go
func main() {
	serviceName := constants.GREETER_SERVICE
	cfg := config.GetConfig()

	grpcOps := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			// keep it last in the interceptor chain
			rpclog.UnaryServerInterceptor(),
		)),
	}

	if cfg.Features.Tls.Enabled {
		tlsConf, err := tls.NewTLSConfig(cfg.Features.Tls.CertFile, cfg.Features.Tls.KeyFile, cfg.Features.Tls.CaFile, cfg.Features.Tls.ServerName)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create cert")
		}
		serverCert := credentials.NewTLS(tlsConf)
		grpcOps = append(grpcOps, grpc.Creds(serverCert))
	}

	srv := service.NewService(
		service.Name(serviceName),
		service.Version(cfg.Services.Greeter.Version),
		service.WithGrpcEndpoint(cfg.Services.Greeter.Endpoint),
		service.WithGrpcOptions(grpcOps...),
		// service.WithBrokerOptions(...),
	)
	// create a gRPC server object
	grpcServer := srv.Server()

	// create a server instance
	greeterHandler := handler.NewGreeterHandler()

	// attach the Greeter service to the server
	greeterv1.RegisterGreeterServiceServer(grpcServer, greeterHandler)

	// start the server
	log.Info().Msg(config.GetBuildInfo())
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
```

## 🔗 Credits
https://github.com/infobloxopen/atlas-app-toolkit/tree/master/server
https://github.com/spencer-p/moroncloudevents/tree/master
