package sub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	broker "github.com/xmlking/toolkit/broker/pubsub"
	terrors "github.com/xmlking/toolkit/errors"
	"github.com/xmlking/toolkit/errors/categories"
	"github.com/xmlking/toolkit/errors/codes"
)

type MySub struct {
	outPublisher broker.Publisher
	dlqPublisher broker.Publisher
}

func NewMySub() *MySub {
	outPublisher, _ := broker.NewPublisher("toolkit-out-dev",
		// broker.WithPublishSettings(pubSettings),
		broker.PublishAsync(true),
	)
	dlqPublisher, _ := broker.NewPublisher("toolkit-in-dead-dev",
		// broker.WithPublishSettings(pubSettings),
		broker.PublishAsync(true),
	)
	return &MySub{
		outPublisher,
        dlqPublisher,
	}
}

func (s *MySub) Handle(ctx context.Context, msg *pubsub.Message) {
	//md, _ := metadata.FromContext(ctx)
	//log.Info().Interface("md", md).Send()
	ID := msg.ID

	// recover from panic
	defer func() {
		if r := recover(); r != nil {

			if err, ok := r.(error); ok {
				log.Error().Err(err).Str("ID", ID).Msg("panic")
				log.Error().Err(err).Str("ID", ID).Msgf("panic details: %+v", err) // for detailed logging
			} else {
				log.Error().Err(err).Interface("cause", r).Str("ID", ID).Msg("panic")
			}

			// Nack will result in the Message being redelivered more quickly then if it were allowed to expire.
			// so we should NOT Nack() and let it expire???
			msg.Nack()

		} else {
			msg.Ack()
		}
	}()

	// to handle server cancellation, We have two options:
	// Option 1: Check if ctx is done, then --> cancel ongoin background work(gRPC, DB, Ongoing PubSub calls) --> msg.Nack()
	// Option 2: Should wait for ongoing background job gracefully finish, then Broker automatically shutdown when all running goroutines finish.

	// Options 1:
	//for {
	//	select {
	//	case <-ctx.Done():
	//		log.Info().Err(err).Msg("Contextended by broker. nacking")
	//		msg.Nack()
	//		return
	//	default:
	//		// keep working
	//	}
	//}

	// Option 2: do processing message work

	log.Info().Interface("event.Message.ID", msg.ID).Send()
	log.Info().Interface("event.Message.Attributes", msg.Attributes).Send()
	log.Info().Interface("event.Message", msg).Send()

    if err := doSomeWork(msg.Data); err != nil {
        s.handleError(ctx, msg, err)
    }

}

func doSomeWork(data []byte) error {
    log.Info().Interface("event.Message.Data", data).Msg("DATA:")
    // TODO: when any processing errors occur, rethrow normalized error and let recover() handle Ack/Nack
    // return terrors.WithCodeAndOperation(errors.New("some sys error"), codes.TempUnavailable, "TempUnavailable")
    // return terrors.WithCodeAndOperation(errors.New("some sys error"), codes.SystemTokenExpired, "SystemTokenExpired")
    return terrors.WithCodeAndOperation(errors.New("some sys error"), codes.DataSchemaNotFound, "DataSchemaNotFound")
}

func (s *MySub) handleError(ctx context.Context, msg *pubsub.Message, err error) {
	headers := msg.Attributes
	if codeErr := terrors.ErrorCoder(nil); errors.As(err, &codeErr) {
		headers["error.code"] = codeErr.Code().String()
		headers["error.category"] = codeErr.Category().String()
		headers["error.desc"] = codeErr.Error()
	}

	switch terrors.GetCategory(err) {
	case categories.System:
		if pubErr := s.dlqPublisher.Publish(ctx, &pubsub.Message{Attributes: headers, Data: msg.Data}); pubErr != nil {
			panic(errors.WithMessage(
				terrors.WithCodeAndOperation(
					errors.WithSecondaryError(pubErr, err),
					codes.SystemTokenExpired, // codes.SystemPublishFailed
					"SystemPublishFailed",
				),
				"Failed publish SystemError to DLQ topic",
			))
		}
		log.Info().Str("ID", msg.ID).Err(err).Msgf("Successfully published errored message to DLQ with headers %v", headers)

	case categories.Data:
		if pubErr := s.outPublisher.Publish(ctx, &pubsub.Message{Attributes: headers, Data: msg.Data}); pubErr != nil {
			panic(errors.WithMessage(
				terrors.WithCodeAndOperation(
					errors.WithSecondaryError(pubErr, err),
					codes.SystemTokenExpired, // codes.SystemPublishFailed
					"SystemPublishFailed",
				),
				"Failed publish DataError to OUTPUT topic",
			))
		}
		log.Info().Str("ID", msg.ID).Err(err).Msgf("Successfully published errored message to OUT with headers %v", headers)

	case categories.Temporary:
		panic(errors.WithMessage(err, "Encountered temporary error"))

	default:
		panic(errors.WithMessage(err, "Encountered uncategorized error, please report"))
	}
}
