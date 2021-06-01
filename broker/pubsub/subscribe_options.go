package broker

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// RecoveryHandler is a function that is called when the recovery middleware recovers from a panic.
// The func takes the receive context, message and the return value from recover
// which reports whether the goroutine is panicking.
// Example usages of HandlerFunc could be to log panics or nack/ack messages which cause panics.
type RecoveryHandler func(context.Context, *pubsub.Message, interface{})

type SubscribeOption func(*SubscribeOptions)

// SubscribeOptions TODO support more pubsub.ReceiveSettings settings
type SubscribeOptions struct {
	// pubsub ReceiveSettings
	ReceiveSettings pubsub.ReceiveSettings

	RecoveryHandler RecoveryHandler
}

// WithRecoveryHandler sets the function for recovering from a panic.
func WithRecoveryHandler(r RecoveryHandler) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.RecoveryHandler = r
	}
}

func WithReceiveSettings(receiveSettings pubsub.ReceiveSettings) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.ReceiveSettings = receiveSettings
	}
}
