package confy

import (
	"context"
	"io/fs"
)

type Option func(*config)

type config struct {
	// Runtime Environment, e.g., test, development, production
	environment string
	// Environment variables prefix. default: "CONFY_"
	environmentVariablePrefix string
	// enable Debug logs. default: false
	debug bool
	// enable Verbose logs. default: false
	verbose bool
	// disable error logs. default: false
	silent bool
	// FileSystem to load config files from. default: os.DirFS(".")
	fs fs.FS
	// errorOnUnmatchedKeys indicating if an error should be thrown
	// if there are keys in the config file that do not correspond to the config struct
	// In case of json files, this field will be used only when compiled with
	// go 1.10 or later.
	// This field will be ignored when compiled with go versions lower than 1.10.
	errorOnUnmatchedKeys bool

	// Alternative options
	context context.Context
}

// Options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithEnvironment set runtime env
func WithEnvironment(env string) Option {
	return func(args *config) {
		args.environment = env
	}
}

// WithEnvironmentVariablePrefix set runtime env
func WithEnvironmentVariablePrefix(envPrefix string) Option {
	return func(args *config) {
		args.environmentVariablePrefix = envPrefix
	}
}

// WithDebugMode enable Debug logs.
// default: false
func WithDebugMode() Option {
	return func(args *config) {
		args.debug = true
	}
}

// WithVerboseMode enable Verbose logs.
// default: false
func WithVerboseMode() Option {
	return func(args *config) {
		args.verbose = true
	}
}

// WithSilentMode disable Error logs.
// default: false
func WithSilentMode() Option {
	return func(args *config) {
		args.silent = true
	}
}

// WithFS enables use custom FileSystem to load config files. e.g., embed.FS
// default: os.DirFS(".")
func WithFS(fs fs.FS) Option {
	return func(args *config) {
		args.fs = fs
	}
}

// WithErrorOnUnmatchedKeys sets if an error should be thrown if
// there are keys in the config file that do not correspond to the config struct.
// default: false
func WithErrorOnUnmatchedKeys() Option {
	return func(args *config) {
		args.errorOnUnmatchedKeys = true
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *config) {
		if o.context == nil {
			o.context = context.Background()
		}
		o.context = context.WithValue(o.context, k, v)
	}
}
