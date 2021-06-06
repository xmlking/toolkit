package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/xmlking/toolkit/broker/pubsub"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	broker.DefaultBroker = broker.NewBroker(appCtx, broker.ProjectID("my-project-id"))

	msg := pubsub.Message{
		ID:         uuid.New().String(),
		Data:       []byte("ABC€"),
		Attributes: map[string]string{"sumo": "demo"},
	}

	var err error
	var pub broker.Publisher
	if pub, err = broker.NewPublisher("toolkit-in-dev",
		broker.WithPublishSettings(pubsub.PublishSettings{Timeout: time.Second * 5})); err != nil {
		log.Fatal().Err(err).Msg("Failed to create publisher for topic: toolkit-in-dev")
	}

	if err := pub.Publish(context.Background(), &msg); err != nil {
		log.Error().Err(err).Msg("Failed to publish to topic: toolkit-in-dev")
	}

	stop()
	log.Info().Msg("Shutting down gracefully, press Ctrl+C again to force")
}
