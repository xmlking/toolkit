package broker

import (
	"context"
)

type SubscribeOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOption func(*SubscribeOptions)

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
