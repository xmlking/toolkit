package broker

import (
	"cloud.google.com/go/pubsub"
)

// TODO support more pubsub.PublishSettings settings
type PublishOptions struct {
	// pubsub PublishSettings
	PublishSettings pubsub.PublishSettings
	// publishes msg to the topic asynchronously if set to true.
	// Default false. i.e., publishes synchronously(blocking)
	Async bool
}

type PublishOption func(*PublishOptions)

func PublishAsync(b bool) PublishOption {
	return func(o *PublishOptions) {
		o.Async = b
	}
}

func WithPublishSettings(publishSettings pubsub.PublishSettings) PublishOption {
	return func(o *PublishOptions) {
		o.PublishSettings = publishSettings
	}
}
