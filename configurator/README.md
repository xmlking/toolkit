# Configurator

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
- Support embed config files in Go binaries via [pkger](https://github.com/markbates/pkger)


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
	"github.com/xmlking/configurator"
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
    configurator.Load(&Config, "config.yml")
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
`debug mode` will let you know how `configurator` loaded your configurations, 
like from which file, shell env, `verbose mode` will tell you even more, like those shell environments `configurator` tried to load.

```go
// Enable debug mode or set env `CONFIG_DEBUG_MODE` to true when running your application
configurator.NewConfigurator(configurator.WithDebugMode()).Load(&Config, "config.yaml")

// Enable verbose mode or set env `CONFIG_VERBOSE_MODE` to true when running your application
configurator.New(configurator.WithVerboseMode()).Load(&Config, "config.yaml")

// You can create custom configurator once and reuse to load multiple different configs  
configor := configurator.NewConfigurator(configurator.WithVerboseMode())
configor.Load(&Config2, "config2.yaml")
configor.Load(&Config3, "config3.yaml")
```

## Load

# Advanced Usage

* Load mutiple configurations

```go
// Earlier configurations have higher priority
configurator.Load(&Config, "application.yml", "database.json")
```

* Return error on unmatched keys

Return an error on finding keys in the config file that do not match any fields in the config struct.
In the example below, an error will be returned if config.toml contains keys that do not match any fields in the ConfigStruct struct.
If ErrorOnUnmatchedKeys is not set, it defaults to false.

Note that for json files, setting ErrorOnUnmatchedKeys to true will have an effect only if using go 1.10 or later.

```go
err := configurator.NewConfigurator(configurator.WithErrorOnUnmatchedKeys()).Load(&ConfigStruct, "config.toml")
```

* Load configuration by environment

Use `CONFIG_ENV` to set environment, if `CONFIG_ENV` not set, environment will be `development` by default, and it will be `test` when running tests with `go test`

```go
// config.go
configurator.Load(&Config, "config.json")

$ go run config.go
// Will load `config.json`, `config.development.json` if it exists
// `config.development.json` will overwrite `config.json`'s configuration
// You could use this to share same configuration across different environments

$ CONFIG_ENV=production go run config.go
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration

$ go test ./configurator/...
// Will load `config.json`, `config.test.json` if it exists
// `config.test.json` will overwrite `config.json`'s configuration

$ CONFIG_ENV=production go test ./configurator/...
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration
```

```go
// Set environment by config
configurator.NewConfigurator(configurator.WithEnvironment("production")).Load(&Config, "config.json")
```

* Example Configuration

```go
// config.go
configurator.Load(&Config, "config.yml")

$ go run config.go
// Will load `config.example.yml` automatically if `config.yml` not found and print warning message
```

* Load files Via [Pkger](https://github.com/markbates/pkger)

> Enable Pkger or set via env `CONFIG_VERBOSE_MODE` to true to use Pkger for loading files

```go
// config.go
configurator.NewConfigurator(configurator.WithPkger()).Load(&Config, "/config/config.json")
# or set via Environment Variable 
$ CONFIG_USE_PKGER=true  go run config.go
```

* Load From Shell Environment

Environment variable names are structured like this:

```
[PREFIX][SEP][MY][SEP][FIELD][SEP][NAME]
```

Struct field names will be converted to **UpperSnakeCase**
```go
$ CONFIG_APP_NAME="hello world" CONFIG_DB_NAME="hello world" go run config.go
// Load configuration from shell environment, it's name is {{prefix}}_FieldName
```

```go
// You could overwrite the prefix with environment CONFIG_ENV_PREFIX, for example:
$ CONFIG_ENV_PREFIX="WEB" WEB_APP_NAME="hello world" WEB_DB_NAME="hello world" go run config.go

// Set prefix by config
configurator.NewConfigurator(configurator.WithEnvironmentVariablePrefix("WEB")).Load(&Config, "config.json")
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

With the `anonymous:"true"` tag specified, the environment variable for the `Description` field is `CONFIG_DESCRIPTION`.
Without the `anonymous:"true"`tag specified, then environment variable would include the embedded struct name and be `CONFIG_DETAILS_DESCRIPTION`.

* With flags

```go
func main() {
	config := flag.String("file", "config.yml", "configuration file")
	flag.StringVar(&Config.APPName, "name", "", "app name")
	flag.StringVar(&Config.DB.Name, "db-name", "", "database name")
	flag.StringVar(&Config.DB.User, "db-user", "root", "database user")
	flag.Parse()

	os.Setenv("CONFIG_ENV_PREFIX", "-")
    configurator.Load(&Config, *config)
	// configurator.Load(&Config) // only load configurations from shell env & flag
}
```

## Gotchas
- Defaults not initialized for `Map` type fields
- Overlaying (merging) not working for `Map` type fields

## TODO
- use [mergo](https://github.com/imdario/mergo) to merge to fix Overlaying `Map` type fields
- Adopt Environment Variables Expanding from [sherifabdlnaby/configuro](https://github.com/sherifabdlnaby/configuro)
- Fully support complex structures involving maps, arrays and slices  [EnvConfig](https://github.com/jlevesy/envconfig)
