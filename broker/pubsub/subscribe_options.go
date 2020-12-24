package broker

import (
	"cloud.google.com/go/pubsub"
	"context"
)

// TODO support more pubsub.ReceiveSettings settings
type SubscribeOptions struct {
	// pubsub ReceiveSettings
	ReceiveSettings pubsub.ReceiveSettings

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

func WithReceiveSettings(receiveSettings pubsub.ReceiveSettings) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.ReceiveSettings = receiveSettings
	}
}
