# Confy

Opinionated configuration loading library for Containerized and 12-Factor compliant applications.

Read configurations from Configuration Files and/or Environment Variables.

This is based on [jinzhu/configor's](https://github.com/jinzhu/configor) and [sherifabdlnaby/configuro's](https://github.com/sherifabdlnaby/configuro) work, with some bug fixes and enhancements. 

## Features

- Strongly typed config with tags
- Reflection based config validation, for syntax amd examples refer [validator](https://github.com/go-playground/validator#examples)
    - Required fields
    - Optional fields
    - Enum fields
    - Min, Max, email, phone etc
- Setting defaults for fields not in the config files. for syntax amd examples refer [creasty's defaults](https://github.com/creasty/defaults)
- Config Sources
    - YAML files
    - Environment Variables
    - [ ] Environment Variables Expanding
    - Command line flags
    - [ ] Kubernetes ConfigMaps
    - Merge multiple config sources, recursively with [Mergo](https://github.com/imdario/mergo)
    - Detect Runtime Environment (test, development, production), auto merge overlay files
- Dynamic Configuration Management (Hot Reconfiguration)
    - Remote config push
    - Externalized configuration
    - Live component reloading / zero-downtime    
    - [ ] Observe Config [Changes](https://play.golang.org/p/41ygGZ-QaB https://gist.github.com/patrickmn/1549985)
- Support embed config files in Go binaries via Go 1.16 [embed](https://golangtutorial.dev/tips/embed-files-in-go/)


**all struct fields must be public**

```golang
type Item struct {
    Name int `yaml:"full_name,omitempty"`
    Age int  `yaml:",omitempty"` //  Removing Empty JSON Fields
    City string `yaml:",omitempty"`
    TLS      bool   `default:"true" yaml:",omitempty"` // Use default when Empty
    Password string `yaml:"-"` // Ignoring Private Fields
    Name     string `validate:"-"`
    Title    string `validate:"alphanum,required"`
    AuthorIP string `validate:"ipv4"`
    Email    string `validate:"email"`
}
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/xmlking/confy"
)

var Config = struct {
	APPName string `default:"app name" yaml:",omitempty"`

	DB struct {
		Name     string
		User     string `default:"root" yaml:",omitempty"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306" yaml:",omitempty"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
    confy.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}
```

With configuration file *config.yml*:

```yaml
appname: test

db:
    name:     test
    user:     test
    password: test
    port:     1234

contacts:
- name: i test
  email: test@test.com
```

## Debug Mode & Verbose Mode

Debug/Verbose mode is helpful when debuging your application, 
`debug mode` will let you know how `confy` loaded your configurations, 
like from which file, shell env, `verbose mode` will tell you even more, like those shell environments `confy` tried to load.

```go
// Enable debug mode or set env `CONFY_DEBUG_MODE` to true when running your application
confy.DefaultConfy = confy.NewConfy(confy.WithFS(os.DirFS(dir)))
confy.Load(&Config, "config.yaml")

// Enable verbose mode or set env `CONFY_VERBOSE_MODE` to true when running your application
confy.New(confy.WithVerboseMode()).Load(&Config, "config.yaml")

// You can create custom confy once and reuse to load multiple different configs  
configor := confy.NewConfy(confy.WithVerboseMode(), confy.WithFS(os.DirFS(".")))
configor.Load(&Config2, "config2.yaml")
configor.Load(&Config3, "config3.yaml")
```

## Load

# Advanced Usage

* Load mutiple configurations

```go
// Earlier configurations have higher priority
confy.Load(&Config, "application.yml", "database.json")
```

* Return error on unmatched keys

Return an error on finding keys in the config file that do not match any fields in the config struct.
In the example below, an error will be returned if config.toml contains keys that do not match any fields in the ConfigStruct struct.
If ErrorOnUnmatchedKeys is not set, it defaults to false.

Note that for json files, setting ErrorOnUnmatchedKeys to true will have an effect only if using go 1.10 or later.

```go
err := confy.NewConfy(confy.WithErrorOnUnmatchedKeys(), confy.WithFS(os.DirFS("."))).Load(&ConfigStruct, "config.toml")
```

* Load configuration by environment

Use `CONFY_ENV` to set environment, if `CONFY_ENV` not set, environment will be `development` by default, and it will be `test` when running tests with `go test`

```go
// config.go
confy.Load(&Config, "config.json")

$ go run config.go
// Will load `config.json`, `config.development.json` if it exists
// `config.development.json` will overwrite `config.json`'s configuration
// You could use this to share same configuration across different environments

$ CONFY_ENV=production go run config.go
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration

$ go test ./confy/...
// Will load `config.json`, `config.test.json` if it exists
// `config.test.json` will overwrite `config.json`'s configuration

$ CONFY_ENV=production go test ./confy/...
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration
```

```go
// Set environment by config
confy.NewConfy(confy.WithEnvironment("production"), confy.WithFS(os.DirFS("."))).Load(&Config, "config.json")
```

* Example Configuration

```go
// config.go
confy.Load(&Config, "config.yml")

$ go run config.go
// Will load `config.example.yml` automatically if `config.yml` not found and print warning message
```

* Load files Via Go 1.16 [FileSystem](https://go.googlesource.com/proposal/+/master/design/draft-iofs.md)


* Load From Shell Environment

Environment variable names are structured like this:

```
[PREFIX][SEP][MY][SEP][FIELD][SEP][NAME]
```

Struct field names will be converted to **UpperSnakeCase**
```go
$ CONFY_APP_NAME="hello world" CONFY_DB_NAME="hello world" go run config.go
// Load configuration from shell environment, it's name is {{prefix}}_FieldName
```

```go
// You could overwrite the prefix with environment CONFY_ENV_PREFIX, for example:
$ CONFY_ENV_PREFIX="WEB" WEB_APP_NAME="hello world" WEB_DB_NAME="hello world" go run config.go

// Set prefix by config
confy.NewConfy(confy.WithEnvironmentVariablePrefix("WEB"), confy.WithFS(os.DirFS("."))).Load(&Config, "config.json")
```

* Anonymous Struct

Add the `anonymous:"true"` tag to an anonymous, embedded struct to NOT include the struct name in the environment
variable of any contained fields.  For example:

```go
type Details struct {
	Description string
}

type Config struct {
	Details `anonymous:"true"`
}
```

With the `anonymous:"true"` tag specified, the environment variable for the `Description` field is `CONFY_DESCRIPTION`.
Without the `anonymous:"true"`tag specified, then environment variable would include the embedded struct name and be `CONFY_DETAILS_DESCRIPTION`.

* With flags

```go
func main() {
	config := flag.String("file", "config.yml", "configuration file")
	flag.StringVar(&Config.APPName, "name", "", "app name")
	flag.StringVar(&Config.DB.Name, "db-name", "", "database name")
	flag.StringVar(&Config.DB.User, "db-user", "root", "database user")
	flag.Parse()

	os.Setenv("CONFY_ENV_PREFIX", "-")
    confy.Load(&Config, *config)
	// confy.Load(&Config) // only load configurations from shell env & flag
}
```

## Gotchas
- Defaults not initialized for `Map` type fields
- Overlaying (merging) not working for `Map` type fields

## TODO
- use [mergo](https://github.com/imdario/mergo) to merge to fix Overlaying `Map` type fields
- check [conflate's mergo](https://github.com/miracl/conflate/blob/master/merge.go)
- Adopt Environment Variables Expanding from [sherifabdlnaby/configuro](https://github.com/sherifabdlnaby/configuro)
- Fully support complex structures involving maps, arrays and slices  [EnvConfig](https://github.com/jlevesy/envconfig)
