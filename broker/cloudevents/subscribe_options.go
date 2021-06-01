package broker

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
)

// RecoveryHandler is a function that is called when the recovery middleware recovers from a panic.
// The func takes the receive context, message and the return value from recover
// which reports whether the goroutine is panicking.
// Example usages of HandlerFunc could be to log panics or nack/ack messages which cause panics.
type RecoveryHandler func(context.Context, event.Event, interface{})

type SubscribeOptions struct {
	RecoveryHandler RecoveryHandler
}

type SubscribeOption func(*SubscribeOptions)

// WithRecoveryHandler sets the function for recovering from a panic.
func WithRecoveryHandler(r RecoveryHandler) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.RecoveryHandler = r
	}
}
