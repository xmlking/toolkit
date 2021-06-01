package broker

import (
	"context"
)

type Option func(*Options)

type Options struct {
	Context context.Context
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
