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
		tlsConf, err := tls.NewTLSConfig(cfg.Features.Tls.CertFile, cfg.Features.Tls.KeyFile, cfg.Features.Tls.CaFile, cfg.Features.Tls.ServerName, cfg.Features.Tls.Password)
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
		service.WithGrpcOptions(grpcOps...), // optional
        // optionally add broker
		service.WithBrokerOptions(
			broker.ProjectID(cfg.Pubsub.ProjectID),
			// broker.ClientOption(option.WithCredentialsFile("GOOGLE_APPLICATION_CREDENTIALS_FILE_PATH")),
		),
	)
	// create a gRPC server object
	grpcServer := srv.Server()

	// create a server instance
	greeterHandler := handler.NewGreeterHandler()

	// attach the Greeter service to the server
	greeterv1.RegisterGreeterServiceServer(grpcServer, greeterHandler)

    testSubscriber := subscriber.testSubscriber()
    
    // optionally add subscribe for broker
	log.Info().Interface("ReceiveSettings", cfg.Pubsub.ReceiveSettings).Send()
	if err := bkr.Subscribe(
		cfg.Pubsub.InputSubscription,
		accountSubscriber.Handle,
		broker.WithReceiveSettings(pubsub.ReceiveSettings(cfg.Pubsub.ReceiveSettings)),
	); err != nil {
		log.Error().Err(err).Msgf("Failed subscribing to Topic: %s", cfg.Sources.Acro.InputTopic)
	}

	// start the server
	log.Info().Msg(config.GetBuildInfo())
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
```

## Run 

### PubSub

Source the script needed for next steps

```bash
. ./scripts/pubsub_functions.sh
```

#### Start PubSub

Start emulator via gcloud cli

```bash
gcps
```

As alternative, you can also start emulator via docker

```bash
docker-compose up pub-sub-emulator
```

#### Setup PubSub

```bash
gcpg
# or 
gcpg tooklit
# or 
gcpg tooklit dev
```

#### Tail logs

```bash
# when using gcloud cli to start emulator
gcpl
```

#### Stop PubSub

```bash
# when using gcloud cli to start emulator
gcpk
# or if you are using docker-compose
docker-compose up down
```

## Development

### Build

```bask
make upgrade_deps
make lint
make format
```

### Test

```bask
make test-unit
```

## 🔗 Credits
https://github.com/infobloxopen/atlas-app-toolkit/tree/master/server
https://github.com/spencer-p/moroncloudevents/tree/master
