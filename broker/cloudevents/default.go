package broker

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	eventing "github.com/xmlking/toolkit/broker/cloudevents/internal"
)

const (
	DefaultName = "mkit.broker.default"
)

type ceBroker struct {
	options Options
	subs    []*ceSubscriber
	pubs    []Publisher
}

type cePublisher struct {
	options  PublishOptions
	ceClient cloudevents.Client
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
		options:  options,
		ceClient: eventing.NewSourceClient(topic),
	}

	b.pubs = append(b.pubs, pub)
	return pub, nil
}

type ceSubscriber struct {
	target   string
	ceClient cloudevents.Client
	options  SubscribeOptions
	hdlr     Handler
}

func (s *ceSubscriber) start(ctx context.Context) (err error) {
	log.Info().Str("component", "cloudevents").Msgf("Subscribing to: %s", s.target)
	if err = s.ceClient.StartReceiver(ctx, s.hdlr); err == nil {
		log.Info().Str("component", "cloudevents").Msgf("Stopped Subscriber Gracefully: %s", s.target)
	}
	return
}

func (b *ceBroker) AddSubscriber(subscription string, hdlr Handler, opts ...SubscribeOption) (err error) {
	options := SubscribeOptions{}

	for _, o := range opts {
		o(&options)
	}

	sub := &ceSubscriber{
		target:   subscription,
		ceClient: eventing.NewSinkClient(subscription),
		options:  options,
		hdlr:     hdlr,
	}

	b.subs = append(b.subs, sub)
	return nil
}

func (b *ceBroker) Start() error {
	ctx := b.options.Context
	g, ctxx := errgroup.WithContext(ctx)

	// start subscribers in the background.
	// when context cancelled, they exit without error.
	for _, sub := range b.subs {
		g.Go(func() error {
			return sub.start(ctxx)
		})
	}

	g.Go(func() (err error) {
		// listen for the interrupt signal
		<-ctx.Done()

		// log situation
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Debug().Str("component", "cloudevents").Msg("Context timeout exceeded")
		case context.Canceled:
			log.Debug().Str("component", "cloudevents").Msg("Context cancelled by interrupt signal")
		}

		log.Info().Str("component", "cloudevents").Msgf("Closing cloudevents client...")

		return nil
	})

	// Wait for all tasks to be finished or return if error occur at any task.
	return g.Wait()
}

// NewBroker creates a new cloudevents broker
func newBroker(ctx context.Context, opts ...Option) Broker {
	// Default Options
	options := Options{
		Name:    DefaultName,
		Context: ctx,
	}

	for _, o := range opts {
		o(&options)
	}

	b := ceBroker{options: options}
	return &b
}
