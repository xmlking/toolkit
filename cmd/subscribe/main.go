package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/xmlking/toolkit/broker/pubsub"
	"github.com/xmlking/toolkit/util/signals"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
)

func main() {
	broker.DefaultBroker = broker.NewBroker(signals.NewContext())

	myHandler := func(ctx context.Context, msg *pubsub.Message) {
		//md, _ := metadata.FromContext(ctx)
		//log.Info().Interface("md", md).Send()
		log.Info().Interface("event.Message.ID", msg.ID).Send()
		log.Info().Interface("event.Message.Attributes", msg.Attributes).Send()
		log.Info().Interface("event.Message.Data", msg.Data).Send()

		log.Info().Interface("event.Message", msg).Send()
		msg.Ack() // or msg.Nack()
	}

	if err := broker.Subscribe("toolkit-in-dev", myHandler, broker.WithSubscriptionID("toolkit-in-dev")); err != nil {
		log.Error().Err(err).Msg("Failed subscribing to Topic: ingestion-in-dev")
	}

	broker.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Info().Msg("Got to Go...")
	// close all subs and then connection.
	if err := broker.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("Unexpected disconnect error")
	}
}
