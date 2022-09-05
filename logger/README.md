# Logger

This logger basically configure **zerolog** so that you can log via `github.com/rs/zerolog/log`

## Usage

### Initialization

Import `shared/logger/auto` package. It will be *self-initialized*.

```golang
import  "github.com/xmlking/toolkit/logger/auto"
```

Other option to initialize logger is to set DefaultLogger your self. this will give more control to developer. 

```golang
logger.DefaultLogger =logger.NewLogger()
logger.DefaultLogger =logger.NewLogger(logger.WithLevel(zerolog.DebugLevel), logger.WithFormat(logger.PRETTY))
logger.DefaultLogger =logger.NewLogger(logger.WithLevel(zerolog.DebugLevel), logger.WithFormat(logger.PRETTY), logger.EnableGrpcLog(true))
```

Once logger is initialized, then you can use standard `github.com/rs/zerolog/log` package's helper methods to log in your code.

### Environment Variables

Your can set **Logger** config via Environment Variables

**grpc** logs are disabled by default. you can enable via `CONFY_LOG_GRPC`

> **grpc** internal logs also adopt `CONFY_LOG_LEVEL` and `CONFY_LOG_FORMAT`

> No need to set `GRPC_GO_LOG_SEVERITY_LEVEL` and `GRPC_GO_LOG_VERBOSITY_LEVEL`

```shell
CONFY_LOG_LEVEL=<trace,debug,info,warn,error,fatal,panic>
CONFY_LOG_FORMAT=<pretty/json/gcp>
CONFY_LOG_GRPC=true
CONFY_LOG_FILE=app1.log
```

Aside from logging in JSON, you can also configure Zerolog to output binary logs encoded in [CBOR](https://cbor.io/) format. You can enable it by using the `binary_log` build tag while compiling your application:

```shell
go build -tags binary_log -o build ./service/engine/...
```

You can decode this binary log entry to JSON with any CBOR decoder, such as [csd](https://github.com/toravir/csd/)
```shell
CONFY_LOG_FORMAT=json ./build/engine  2> >(csd)
```

## Test
```shell
CONFY_LOG_LEVEL=info  CONFY_LOG_FORMAT=json go test github.com/xmlking/toolkit/logger  -count=1
```


## Reference 
- A Complete Guide to Logging in Go with [Zerolog](https://betterstack.com/community/guides/logging/zerolog/)
