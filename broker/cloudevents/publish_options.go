package broker

import (
	"context"
)

type PublishOptions struct {
	// publishes msg to the topic asynchronously if set to true.
	// Default false. i.e., publishes synchronously(blocking)
	Async bool
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type PublishOption func(*PublishOptions)

// PublishContext set context
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

func PublishAsync(b bool) PublishOption {
	return func(o *PublishOptions) {
		o.Async = b
	}
}
