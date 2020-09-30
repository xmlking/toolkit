# Middleware

1. rpclog
1. duration
1. tags
1. [validator](github.com/grpc-ecosystem/go-grpc-middleware/validator)
1. [errors](https://github.com/cockroachdb/errors/blob/master/grpc/main_test.go)
1. [opentelemetry](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/google.golang.org/grpc)
1. [opentracing](https://github.com/grpc-ecosystem/grpc-opentracing)
1. [go-concurrency-limits](https://github.com/platinummonkey/go-concurrency-limits)



## Get

```bash
go get github.com/xmlking/toolkit
```

## Usage

Interceptors will be executed **from left to right**: e.g., logging, monitoring and auth.
```go
grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(loggingUnary, monitoringUnary, authUnary),)
```

Add interceptors in following order
1. Around interceptors - from outer to inner — e.g., duration,  retry  
2. Before interceptors - rate-limit, auth, validation , tagging 
3. After interceptors - rpclog, translog, recovery

```go
import (
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
    "github.com/xmlking/toolkit/middleware/rpclog"
)

server := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        // Execution is done in left-to-right order
        grpc_validator.UnaryServerInterceptor(),
        // keep it last in the interceptor chain
        rpclog.UnaryServerInterceptor(rpclog.WithExcludeMethods("/grpc.health.v1.Health/Check", "/api.MyService/*")),
    )),
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        // keep it last in the interceptor chain
        rpclog.StreamServerInterceptor()
    )),
)
```

#### concurrency-limits

```go
import (
    gclGrpc "github.com/platnummonkey/go-concurrency-limits/grpc"
)

// setup grpc server with this option
serverOption := grpc.UnaryInterceptor(
    gclGrpc.UnaryServerInterceptor(
        gclGrpc.WithLimiter(...),
        gclGrpc.WithServerResponseTypeClassifier(..),
    ),
)

// setup grpc client with this option
dialOption := grpc.WithUnaryInterceptor(
    gclGrpc.UnaryClientInterceptor(
        gclGrpc.WithLimiter(...),
        gclGrpc.WithClientResponseTypeClassifier(...),
    ),
)
```
