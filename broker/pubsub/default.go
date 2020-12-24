package broker

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pubsubBroker struct {
	client  *pubsub.Client
	options Options
	subs    []Subscriber
	pubs    []Publisher
}

type pubsubPublisher struct {
	options PublishOptions
	topic   *pubsub.Topic
}

func (p *pubsubPublisher) Topic() string {
	return p.topic.String()
}

// Stop should be called once
func (p *pubsubPublisher) Stop() {
	log.Info().Msgf("Stopping Publisher: %s", p.Topic())
	p.topic.Stop()
	log.Info().Msgf("Stopped Publisher Gracefully: %s", p.Topic())
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

func (s *pubsubSubscriber) Start() {
	for {
		log.Info().Msgf("Subscribing: %s", s.sub)
		if err := s.sub.Receive(s.options.Context, s.hdlr); err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Canceled:
					log.Info().Msgf("Subscriber(%s) Canceled. Stopping Subscription", s.sub)
					break
				default:
					log.Error().Err(err).Msgf("Subscriber(%s) Error. Retrying...", s.sub)
				}
			} else {
				log.Error().Err(err).Msgf("Subscriber(%s) Error Unknown. Retrying...", s.sub)
			}
			// got error while subscribing to topic. lets retry after 1 sec.
			time.Sleep(time.Second)
			continue
		} else {
			// ctx is done. gracefully exiting the loop
			break
		}
	}
	close(s.done)
}

// Stop should be called once
func (s *pubsubSubscriber) Stop() {
	log.Info().Msgf("Stopping Subscriber: %s", s.sub)
	for {
		select {
		case <-s.done:
			log.Info().Msgf("Stopped Subscriber Gracefully: %s", s.sub)
			return
		}
	}
}

// Shutdown shuts down all subscribers gracefully and then close the connection
func (b *pubsubBroker) Shutdown() (err error) {
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
	err = b.client.Close()
	// Hint: when using pubsub emulator, you receive this error, which you can safely ignore.
	// Live pubsub server will throw this error.
	if err != nil && strings.Contains(err.Error(), "the client connection is closing") {
		err = nil
	}
	return
}

func (b *pubsubBroker) Options() Options {
	return b.options
}

func (b *pubsubBroker) NewPublisher(topic string, opts ...PublishOption) (pub Publisher, err error) {
	t := b.client.Topic(topic)

	options := PublishOptions{
		Async:   false,
		Context: b.options.Context,
	}

	for _, o := range opts {
		o(&options)
	}

	var exists bool
	exists, err = t.Exists(options.Context)
	if err != nil {
		return
	}
	if !exists {
		err = errors.Errorf("Doesn't exist Topic: %s", t)
		return
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

	pub = &pubsubPublisher{
		topic: t,
	}
	// keep track of pubs
	b.pubs = append(b.pubs, pub)

	return
}

// Subscribe registers a subscription to the given topic against the google pubsub api
func (b *pubsubBroker) Subscribe(subscription string, h Handler, opts ...SubscribeOption) (err error) {
	options := SubscribeOptions{
		Context: b.options.Context,
	}

	for _, o := range opts {
		o(&options)
	}

	sub := b.client.Subscription(subscription)
	exists, err := sub.Exists(context.TODO()) // TODO should we use context.Background()  ?
	if err != nil {
		return err
	}
	if !exists {
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

	subscriber := &pubsubSubscriber{
		options: options,
		done:    make(chan struct{}),
		sub:     sub,
		hdlr:    h,
	}

	// keep track of subs
	b.subs = append(b.subs, subscriber)

	return nil
}

// Start should be called once
func (b *pubsubBroker) Start() error {
	for _, sub := range b.subs {
		// TODO: should we capture start error and return? via errCh <- ?
		go sub.Start()
	}
	return nil
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
