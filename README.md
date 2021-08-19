# toolkit

## Features

- [x] Config
- [x] Logging
- [x] Broker
- [x] Errors
- [x] Server
- [x] Middleware
- [x] Telemetry
- [ ] Auth
- [ ] Cache
- [ ] Crypto

## Usage

```go
func main() {
    serviceName := constants.PLAY_SERVICE
    cfg := config.GetConfig()
    efs := config.GetFileSystem()
    
    appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
    defer stop()
    
    g, ctx := errgroup.WithContext(appCtx)
    
    // Register kuberesolver to grpc.
    // This line should be before calling registry.NewContainer(cfg)
    if config.IsProduction() {
    kuberesolver.RegisterInCluster()
    }
    
    if cfg.Features.Tracing.Enabled {
    closeFn := tracing.InitTracing(ctx, cfg.Features.Tracing)
    defer closeFn()
    }
    
    if cfg.Features.Metrics.Enabled {
    closeFn := metrics.InitMetrics(ctx, cfg.Features.Metrics)
    defer closeFn()
    }
    
    var unaryInterceptors = []grpc.UnaryServerInterceptor{grpc_validator.UnaryServerInterceptor()}
    var streamInterceptors = []grpc.StreamServerInterceptor{grpc_validator.StreamServerInterceptor()}
    
    if cfg.Features.Tracing.Enabled {
    unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryServerInterceptor())
    streamInterceptors = append(streamInterceptors, otelgrpc.StreamServerInterceptor())
    }
    if cfg.Features.Rpclog.Enabled {
    // keep it last in the interceptor chain
    unaryInterceptors = append(unaryInterceptors, rpclog.UnaryServerInterceptor())
    }
    
    // ServerOption
    grpcOps := []grpc.ServerOption{
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
    }
    
    if cfg.Features.TLS.Enabled {
    tlsConf, err := tls.NewTLSConfig(efs, cfg.Features.TLS.CertFile, cfg.Features.TLS.KeyFile, cfg.Features.TLS.CaFile, cfg.Features.TLS.ServerName, cfg.Features.TLS.Password)
    if err != nil {
    log.Fatal().Err(err).Msg("failed to create cert")
    }
    serverCert := credentials.NewTLS(tlsConf)
    grpcOps = append(grpcOps, grpc.Creds(serverCert))
    }
    
    listener, err := endpoint.GetListener(cfg.Services.Play.Endpoint)
    if err != nil {
    log.Fatal().Stack().Err(err).Msg("error creating listener")
    }
    srv := server.NewServer(appCtx, server.ServerName(serviceName), server.WithListener(listener), server.WithServerOptions(grpcOps...))
    
    gSrv := srv.Server()
    
    greeterHandler := handler.NewGreeterHandler()
    // attach the Greeter service to the server
    greeterv1.RegisterGreeterServiceServer(gSrv, greeterHandler)
    
    // Start broker/gRPC daemon services
    log.Info().Msg(config.GetBuildInfo())
    log.Info().Msgf("Server(%s) starting at: %s, secure: %t, pid: %d", serviceName, listener.Addr(), cfg.Features.TLS.Enabled, os.Getpid())
    
    g.Go(func() error {
    return srv.Start()
    })
    
    go func() {
    if err := g.Wait(); err != nil {
    log.Fatal().Stack().Err(err).Msgf("Unexpected error for service: %s", cfg.Services.Emailer.Endpoint)
    }
    log.Info().Msg("Goodbye.....")
    os.Exit(0)
    }()
    
    // Listen for the interrupt signal.
    <-appCtx.Done()
    
    // notify user of shutdown
    switch ctx.Err() {
    case context.DeadlineExceeded:
    log.Info().Str("cause", "timeout").Msg("Shutting down gracefully, press Ctrl+C again to force")
    case context.Canceled:
    log.Info().Str("cause", "interrupt").Msg("Shutting down gracefully, press Ctrl+C again to force")
    }
    
    // Restore default behavior on the interrupt signal.
    stop()
    
    // Perform application shutdown with a maximum timeout of 1 minute.
    timeoutCtx, cancel := context.WithTimeout(context.Background(), constants.DefaultShutdownTimeout)
    defer cancel()
    
    // force termination after shutdown timeout
    <-timeoutCtx.Done()
    log.Error().Msg("Shutdown grace period elapsed. force exit")
    // force stop any daemon services here:
    srv.Stop()
    os.Exit(1)
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
task mod:download
task go:lint
task go:format
```

### Test

```bask
task go:test
```

### Release

Replace **vx.y.z** with version you try to tag. e.g., **v0.2.3**

1. Start release

    ```bask
    git switch develop
    git flow release start vx.y.z
    ```

2. Update files

    Update  all **go.mod** files that have reference to `github.com/xmlking/toolkit v0.2.3` -> `github.com/xmlking/toolkit vx.y.z`. e.g., 

    ```
    broker/cloudevents/go.mod
    cmd/publish/go.mod
    cmd/subscribe/go.mod
    ```

3. update deps

    ```bask
    task mod:outdated
    # then upgrade recomended versions in each go.mod files
    task mod:sync
    task mod:verify
    ```

4. update changelog

    ```bask
    git-chglog -c .github/chglog/config.yml -o CHANGELOG.md --next-tag vx.y.z
    ```

5. Finish release

    ```bask
    git add .
    git commit -m "build(release): update changelog"
    git flow release finish
    git push origin --all && git push origin --tags # alias for gpoat 
    ```

6. Push tags for all modules

    ```bask
    git switch main
    task mod:release TAG=vx.y.z
    ```

### Release



## 🔗 Credits

https://github.com/infobloxopen/atlas-app-toolkit/tree/master/server
https://github.com/spencer-p/moroncloudevents/tree/master

## Similar Projects

- [Kratos](https://go-kratos.dev/)
    - [Kratos Docs]( https://go-kratos.dev/en/docs/)
    - [Kratos Project Template](https://github.com/go-kratos/kratos-layout)
