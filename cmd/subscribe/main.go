package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/broker/pubsub"
	"golang.org/x/sync/errgroup"
)

const (
	// DefaultShutdownTimeout defines the default timeout given to the service when calling Shutdown.
	DefaultShutdownTimeout = time.Minute * 1
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	g, ctx := errgroup.WithContext(appCtx)

	broker.DefaultBroker = broker.NewBroker(ctx, broker.ProjectID("my-project-id"))

	myHandler := func(ctx context.Context, msg *pubsub.Message) {
		//md, _ := metadata.FromContext(ctx)
		//log.Info().Interface("md", md).Send()
		log.Info().Interface("event.Message.ID", msg.ID).Send()
		log.Info().Interface("event.Message.Attributes", msg.Attributes).Send()
		log.Info().Interface("event.Message.Data", msg.Data).Send()

		log.Info().Interface("event.Message", msg).Send()
		msg.Ack() // or msg.Nack()
	}

	if err := broker.AddSubscriber("toolkit-in-dev", myHandler); err != nil {
		log.Fatal().Err(err).Msg("Failed subscribing to Topic: toolkit-in-dev")
	}

	g.Go(func() error {
		return broker.Start()
	})

	go func() {
		if err := g.Wait(); err != nil {
			log.Fatal().Stack().Err(err).Msg("Unexpected error")
		}
		log.Info().Msg("Goodbye.....")
		os.Exit(0)
	}()

	// Listen for the interrupt signal.
	<-appCtx.Done()

	// notify user of shutdown
	switch ctx.Err() {
	case context.DeadlineExceeded:
		log.Info().Str("cause", "timeout").Msg("Shutting down gracefully, press Ctrl+C again to force")
	case context.Canceled:
		log.Info().Str("cause", "interrupt").Msg("Shutting down gracefully, press Ctrl+C again to force")
	}

	// Restore default behavior on the interrupt signal.
	stop()

	// Perform application shutdown with a maximum timeout of 1 minute.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	// force termination after shutdown timeout
	<-timeoutCtx.Done()
	log.Error().Msg("Shutdown grace period elapsed. force exit")
	os.Exit(1)
}
