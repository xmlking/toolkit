package broker

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Options() Options
	NewPublisher(topic string, opts ...PublishOption) (pub Publisher, err error)
	NewSubscriber(subscription string, hdlr Handler, opts ...SubscribeOption) (sub Subscriber, err error)
	Start() error
	Shutdown() error
}

type Publisher interface {
	Publish(ctx context.Context, msg *pubsub.Message) error
	Stop()
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(context.Context, *pubsub.Message)

// Subscriber is a convenience ~return~ type for the Subscribe method
type Subscriber interface {
	Start()
	Stop()
}

var DefaultBroker Broker

// NewBroker creates and returns a new Broker based on the packages within.
func NewBroker(ctx context.Context, opts ...Option) Broker {
	return newBroker(ctx, opts...)
}

func Start() error {
	return DefaultBroker.Start()
}

func Shutdown() error {
	return DefaultBroker.Shutdown()
}

func NewPublisher(topic string, opts ...PublishOption) (Publisher, error) {
	return DefaultBroker.NewPublisher(topic, opts...)
}

func NewSubscriber(subscription string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	return DefaultBroker.NewSubscriber(subscription, handler, opts...)
}