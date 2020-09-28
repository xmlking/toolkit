# Logger

This logger basically configure **zerolog** so that you can log via `github.com/rs/zerolog/log`

## Usage

Import `shared/logger` package. It will be *self-initialized*. 

```golang
import  "github.com/xmlking/toolkit/logger"
```

Once logger is initialized, then you can use standard `github.com/rs/zerolog/log` package's helper methods to log in your code.



### Environment Variables 

Your can set **Logger** config via Environment Variables

**grpc** logs are disabled by default. you can enable via `CONFIGOR_LOG_GRPC`

> **grpc** internal logs also adopt `CONFIGOR_LOG_LEVEL` and `CONFIGOR_LOG_FORMAT`

> No need to set `GRPC_GO_LOG_SEVERITY_LEVEL` and `GRPC_GO_LOG_VERBOSITY_LEVEL`

```
CONFIGOR_LOG_LEVEL=<trace,debug,info,warn,error,fatal,panic>
CONFIGOR_LOG_FORMAT=<pretty/json/gcp>
CONFIGOR_LOG_GRPC=true
```

## Test
```
CONFIGOR_LOG_LEVEL=info CONFIGOR_LOG_FORMAT=json go test github.com/xmlking/toolkit/logger  -count=1
```
