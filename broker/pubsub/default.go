package broker

import (
	"context"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type pubsubBroker struct {
	client  *pubsub.Client
	options Options
	subs    []*pubsubSubscriber
	pubs    []*pubsubPublisher
}

type pubsubPublisher struct {
	options PublishOptions
	topic   *pubsub.Topic
}

func (p *pubsubPublisher) Topic() string {
	return p.topic.String()
}

// Stop should be called once
func (p *pubsubPublisher) stop() {
	log.Info().Str("component", "pubsub").Msgf("Stopping Publisher: %s", p.Topic())
	// It blocks until all items have been flushed.
	p.topic.Stop()
	log.Info().Str("component", "pubsub").Msgf("Stopped Publisher Gracefully: %s", p.Topic())
}

func (p *pubsubPublisher) Publish(ctx context.Context, msg *pubsub.Message) (err error) {
	pr := p.topic.Publish(ctx, msg)
	if !p.options.Async {
		if _, err = pr.Get(ctx); err != nil {
			log.Error().Err(err).Msgf("Unable to publish to topic: %s", p.topic.String())
		}
	}
	return
}

type pubsubSubscriber struct {
	options SubscribeOptions
	sub     *pubsub.Subscription
	hdlr    Handler
	done    chan struct{}
}

func (s *pubsubSubscriber) start(ctx context.Context) (err error) {
	defer close(s.done)
	log.Info().Str("component", "pubsub").Msgf("Subscribing to: %s", s.sub)
	// If ctx is done, Receive returns nil after all of the outstanding calls to `s.hdlr` have returned
	// and all messages have been acknowledged or have expired.
	if err = s.sub.Receive(ctx, s.hdlr); err == nil {
		log.Info().Str("component", "pubsub").Msgf("Stopped Subscriber Gracefully: %s", s.sub)
	}
	return
}

func (b *pubsubBroker) NewPublisher(topic string, opts ...PublishOption) (Publisher, error) {
	t := b.client.Topic(topic)

	options := PublishOptions{
		Async: false,
	}

	for _, o := range opts {
		o(&options)
	}

	if exists, err := t.Exists(context.Background()); err != nil {
		return nil, err
	} else if !exists {
		err = errors.Errorf("Doesn't exist Topic: %s", t)
		return nil, err
	}

	if options.PublishSettings.DelayThreshold != 0 {
		t.PublishSettings.DelayThreshold = options.PublishSettings.DelayThreshold
	}
	if options.PublishSettings.CountThreshold != 0 {
		t.PublishSettings.CountThreshold = options.PublishSettings.CountThreshold
	}
	if options.PublishSettings.ByteThreshold != 0 {
		t.PublishSettings.ByteThreshold = options.PublishSettings.ByteThreshold
	}
	if options.PublishSettings.NumGoroutines != 0 {
		t.PublishSettings.NumGoroutines = options.PublishSettings.NumGoroutines
	}
	if options.PublishSettings.Timeout != 0 {
		t.PublishSettings.Timeout = options.PublishSettings.Timeout
	}
	if options.PublishSettings.BufferedByteLimit != 0 {
		t.PublishSettings.BufferedByteLimit = options.PublishSettings.BufferedByteLimit
	}

	pub := &pubsubPublisher{
		topic: t,
	}
	// keep track of pubs
	b.pubs = append(b.pubs, pub)

	return pub, nil
}

// AddSubscriber registers a subscription to the given topic against the google pubsub api
func (b *pubsubBroker) AddSubscriber(subscription string, hdlr Handler, opts ...SubscribeOption) error {
	options := SubscribeOptions{}

	for _, o := range opts {
		o(&options)
	}

	sub := b.client.Subscription(subscription)
	if exists, err := sub.Exists(context.Background()); err != nil {
		return err
	} else if !exists {
		return errors.Errorf("Subscription %s doesn't exists", sub)
	}

	if options.ReceiveSettings.MaxOutstandingBytes != 0 {
		sub.ReceiveSettings.MaxOutstandingBytes = options.ReceiveSettings.MaxOutstandingBytes
	}
	if options.ReceiveSettings.MaxOutstandingMessages != 0 {
		sub.ReceiveSettings.MaxOutstandingMessages = options.ReceiveSettings.MaxOutstandingMessages
	}
	if options.ReceiveSettings.NumGoroutines != 0 {
		sub.ReceiveSettings.NumGoroutines = options.ReceiveSettings.NumGoroutines
	}
	if options.ReceiveSettings.MaxExtension != 0 {
		sub.ReceiveSettings.MaxExtension = options.ReceiveSettings.MaxExtension
	}
	if options.ReceiveSettings.MaxExtensionPeriod != 0 {
		sub.ReceiveSettings.MaxExtensionPeriod = options.ReceiveSettings.MaxExtensionPeriod
	}
	if options.ReceiveSettings.Synchronous != false {
		sub.ReceiveSettings.Synchronous = options.ReceiveSettings.Synchronous
	}

	middleware := hdlr
	if rHdlr := options.RecoveryHandler; rHdlr != nil {
		middleware = func(ctx context.Context, msg *pubsub.Message) {
			defer func() {
				if r := recover(); r != nil {
					rHdlr(ctx, msg, r)
				}
			}()

			hdlr(ctx, msg)
		}
	}

	subscriber := &pubsubSubscriber{
		options: options,
		done:    make(chan struct{}),
		sub:     sub,
		hdlr:    middleware,
	}

	// keep track of subs
	b.subs = append(b.subs, subscriber)

	return nil
}

// Start blocking. run as background process.
func (b *pubsubBroker) Start() (err error) {
	ctx := b.options.Context
	g, egCtx := errgroup.WithContext(ctx)

	// start subscribers in the background.
	// when context cancelled, they exit without error.
	for _, sub := range b.subs {
		g.Go(func() error {
			return sub.start(egCtx)
		})
	}

	g.Go(func() (err error) {
		// listen for the interrupt signal
		<-ctx.Done()

		// log situation
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Debug().Str("component", "pubsub").Msg("Context timeout exceeded")
		case context.Canceled:
			log.Debug().Str("component", "pubsub").Msg("Context cancelled by interrupt signal")
		}

		// wait for all subs to stop
		for _, sub := range b.subs {
			log.Info().Str("component", "pubsub").Msgf("Stopping Subscriber: %s", sub.sub)
			<-sub.done
		}

		// then stop all pubs
		for _, pub := range b.pubs {
			pub.stop()
		}

		// then disconnection client.
		log.Info().Str("component", "pubsub").Msgf("Closing pubsub client...")
		err = b.client.Close()

		// Hint: when using pubsub emulator, you receive this error, which you can safely ignore.
		// Live pubsub server will throw this error.
		if err != nil && strings.Contains(err.Error(), "the client connection is closing") {
			err = nil
		}
		return
	})

	// Wait for all tasks to be finished or return if error occur at any task.
	return g.Wait()
}

// NewBroker creates a new google pubsub broker
func newBroker(ctx context.Context, opts ...Option) Broker {
	// Default Options
	options := Options{
		Context: ctx,
	}

	for _, o := range opts {
		o(&options)
	}

	// retrieve project id
	prjID := options.ProjectID

	// if `GOOGLE_CLOUD_PROJECT` is present, it will overwrite programmatically set projectID
	//if envPrjID := os.Getenv("GOOGLE_CLOUD_PROJECT"); len(envPrjID) > 0 {
	//	prjID = envPrjID
	//}

	// create pubsub client
	c, err := pubsub.NewClient(ctx, prjID, options.ClientOptions...)
	if err != nil {
		panic(err.Error())
	}

	return &pubsubBroker{
		client:  c,
		options: options,
	}
}
