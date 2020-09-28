package broker

import (
	"cloud.google.com/go/pubsub"
	"context"
)

// TODO support more pubsub.ReceiveSettings settings
type SubscribeOptions struct {
	// pubsub ReceiveSettings
	ReceiveSettings pubsub.ReceiveSettings
	// Subscribers with the same Subscription ID
	// will create a shared subscription where each
	// receives a subset of messages.
	SubscriptionID string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type SubscribeOption func(*SubscribeOptions)

// WithSubscriptionID sets the SubscriptionID of the topic to share messages on
func WithSubscriptionID(id string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.SubscriptionID = id
	}
}

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
