package broker

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

type Broker interface {
	Options() Options
	NewPublisher(topic string, opts ...PublishOption) (pub Publisher, err error)
	NewSubscriber(subscription string, hdlr Handler, opts ...SubscribeOption) (sub Subscriber, err error)
	Start() error
	Shutdown() error
}

type Handler func(context.Context, event.Event) cloudevents.Result

type Publisher interface {
	Publish(ctx context.Context, event event.Event) error
	Stop()
}

type Subscriber interface {
	Start()
	Stop()
}

// NewBroker creates and returns a new Broker based on the packages within.
func NewBroker(opts ...Option) Broker {
	return newBroker(opts...)
}
