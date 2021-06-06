package broker

import (
	"context"
)

type Option func(*Options)

type Options struct {
	Name    string
	Context context.Context
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Context  appContext to trigger terminate signal
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
