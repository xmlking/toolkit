package broker

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/rs/zerolog/log"

	eventing "github.com/xmlking/toolkit/broker/cloudevents/internal"
)

const (
	DefaultName = "mkit.broker.default"
)

type ceBroker struct {
	opts Options
	subs []Subscriber
	pubs []Publisher
}

type ceSubscriber struct {
	name     string
	target   string
	ceClient cloudevents.Client
	options  SubscribeOptions
	hdlr     Handler
}

func (s *ceSubscriber) Start() {
	log.Info().Msgf("Subscriber (%s) starting at: %s", s.name, s.target)
	if err := s.ceClient.StartReceiver(context.Background(), s.hdlr); err != nil {
		log.Error().Err(err).Send()
	}
}

// Stop should be called once
func (s *ceSubscriber) Stop() {
	log.Info().Msgf("Stopping Subscriber")
	log.Info().Msgf("Stopped Subscriber Gracefully")
}

type cePublisher struct {
	ceClient cloudevents.Client
	options  PublishOptions
}

// Stop should be called once
func (p *cePublisher) Stop() {
	log.Info().Msgf("Stopping Publisher")
	log.Info().Msgf("Stopped Publisher Gracefully")
}

func (p *cePublisher) Publish(ctx context.Context, event event.Event) (err error) {
	return p.ceClient.Send(ctx, event)
}

func (b *ceBroker) NewPublisher(topic string, opts ...PublishOption) (pub Publisher, err error) {
	options := PublishOptions{
		Async: false,
	}

	for _, o := range opts {
		o(&options)
	}

	pub = &cePublisher{
		ceClient: eventing.NewSourceClient(topic),
		options:  options,
	}

	b.pubs = append(b.pubs, pub)
	return pub, nil
}

func (b *ceBroker) NewSubscriber(subscription string, hdlr Handler, opts ...SubscribeOption) (sub Subscriber, err error) {
	options := SubscribeOptions{
		context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	sub = &ceSubscriber{
		name:     b.opts.Name,
		target:   subscription,
		ceClient: eventing.NewSinkClient(subscription),
		options:  options,
		hdlr:     hdlr,
	}

	b.subs = append(b.subs, sub)
	return sub, nil
}

func (b *ceBroker) Start() error {
	for _, sub := range b.subs {
		go sub.Start()
	}
	return nil
}

func (b *ceBroker) Shutdown() (err error) {
	// close all subs
	for _, sub := range b.subs {
		sub.Stop()
	}
	// then close all pubs
	for _, pub := range b.pubs {
		pub.Stop()
	}
	// then disconnection client.
	log.Info().Msgf("Closing pubsub client...")
	return nil
}

func newBroker(opts ...Option) Broker {
	// Default Options
	options := Options{
		Name: DefaultName,
	}
	b := ceBroker{opts: options}
	b.ApplyOptions(opts...)
	return &b
}

func (b *ceBroker) ApplyOptions(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&b.opts)
	}
}

func (b *ceBroker) Options() Options {
	return b.opts
}
