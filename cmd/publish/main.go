package main

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/xmlking/toolkit/broker/pubsub"
)

func main() {
	// broker.DefaultBroker = broker.NewBroker(context.Background(), broker.ProjectID("my-project-id")); // use cfg.pubsub.ProjectID
	broker.DefaultBroker = broker.NewBroker(context.Background())

	msg := pubsub.Message{
		ID:         uuid.New().String(),
		Data:       []byte("ABC€"),
		Attributes: map[string]string{"sumo": "demo"},
	}

	var err error
	var pub broker.Publisher
	if pub, err = broker.NewPublisher("toolkit-in-dev",
		broker.WithPublishSettings(pubsub.PublishSettings{Timeout: time.Second * 5})); err != nil {
		log.Error().Err(err).Send()
	}

	if err := pub.Publish(context.Background(), &msg); err != nil {
		log.Error().Err(err).Send()
	}

	pub.Stop()
}
