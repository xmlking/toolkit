package broker

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	NewPublisher(topic string, opts ...PublishOption) (pub Publisher, err error)
	AddSubscriber(subscription string, hdlr Handler, opts ...SubscribeOption) (err error)
	Start() error
}

type Publisher interface {
	Publish(ctx context.Context, event event.Event) error
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(context.Context, event.Event) cloudevents.Result

var DefaultBroker Broker

// NewBroker creates and returns a new Broker based on the packages within.
func NewBroker(ctx context.Context, opts ...Option) Broker {
	return newBroker(ctx, opts...)
}

func Start() error {
	return DefaultBroker.Start()
}

func NewPublisher(topic string, opts ...PublishOption) (Publisher, error) {
	return DefaultBroker.NewPublisher(topic, opts...)
}

func AddSubscriber(subscription string, handler Handler, opts ...SubscribeOption) error {
	return DefaultBroker.AddSubscriber(subscription, handler, opts...)
}
