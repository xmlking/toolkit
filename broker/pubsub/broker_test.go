package broker_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	broker "github.com/xmlking/toolkit/broker/pubsub"
	"google.golang.org/api/option"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
)

func TestNewBroker(t *testing.T) {
	ctX, cancel := context.WithCancel(context.Background())
	srv := setupFakePubsubAndBroker(ctX, t)
	defer srv.Close()

	myHandler := func(ctx context.Context, msg *pubsub.Message) {
		t.Logf("received msg: %s", msg.Data)
		msg.Ack()
		cancel()
		if got, want := string(msg.Data), "hello"; got != want {
			t.Fatalf(`got "%s" message, want "%s"`, got, want)
			msg.Nack()
		}
	}

	// add subscriber
	if err := broker.AddSubscriber("sumo", myHandler); err != nil {
		t.Fatal(err)
	}

	// start broker
	go broker.Start()

	// create publisher
	publisher, err := broker.NewPublisher("sumo", broker.PublishAsync(false))
	if err != nil {
		t.Fatal(err)
	}

	// publish message
	err = publisher.Publish(context.Background(), &pubsub.Message{Data: []byte("hello")})
	if err != nil {
		t.Fatal(err)
	}

	//srv.Wait()
	ms := srv.Messages()
	t.Logf("msg1: %s", ms[0].Data)
	tps, err := srv.GServer.ListTopics(context.TODO(), &pb.ListTopicsRequest{Project: "projects/pro1"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tps)
}

func TestSubscribeWithRecoveryHandler(t *testing.T) {
	ctX, cancel := context.WithCancel(context.Background())
	srv := setupFakePubsubAndBroker(ctX, t)
	defer srv.Close()

	myHandler := func(ctx context.Context, msg *pubsub.Message) {
		t.Logf("received msg: %s", msg.Data)
		msg.Ack()
		cancel()
		if got, want := string(msg.Data), "hello"; got != want {
			t.Fatalf(`got "%s" message, want "%s"`, got, want)
			msg.Nack()
		}
		panic("sumo")
	}

	recoveryHandler := func(ctx context.Context, msg *pubsub.Message, r interface{}) {
		t.Logf("Recovered from panic: %v,  msg.ID: %s", r, msg.ID)
	}

	// add subscriber
	if err := broker.AddSubscriber("sumo", myHandler, broker.WithRecoveryHandler(recoveryHandler)); err != nil {
		t.Fatal(err)
	}

	// start broker
	go broker.Start()

	// create publisher
	publisher, err := broker.NewPublisher("sumo", broker.PublishAsync(false))
	if err != nil {
		t.Fatal(err)
	}

	// publish message
	err = publisher.Publish(context.Background(), &pubsub.Message{Data: []byte("hello")})
	if err != nil {
		t.Fatal(err)
	}

	//srv.Wait()
	ms := srv.Messages()
	t.Logf("msg1: %s", ms[0].Data)
	tps, err := srv.GServer.ListTopics(context.TODO(), &pb.ListTopicsRequest{Project: "projects/pro1"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tps)
}

// setupFakePubsubAndBroker creates a new fake pubsub server and setup topics and subscriptions
// it also created broker and set as default.
// Note: be sure to close server and broker!
func setupFakePubsubAndBroker(ctx context.Context, t *testing.T) *pstest.Server {
	t.Helper()

	srv := pstest.NewServer()
	// Connect to the server without using TLS.
	conn, err := grpc.DialContext(ctx, srv.Addr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	client, err := pubsub.NewClient(ctx, "pro1", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatal(err)
	}

	ctxl := context.Background()

	// srv.GServer.CreateTopic(ctxl, &pb.Topic{Name:"projects/pro1/topics/sumo"})
	// create topic "projects/pro1/topics/sumo"
	top, err := client.CreateTopic(ctxl, "sumo")
	if err != nil {
		t.Fatal(err)
	}

	// srv.GServer.CreateSubscription(ctxl, &pb.Subscription{Name:"projects/pro1/subscriptions/sumo", Topic: top.String()})
	// create subscription "projects/pro1/subscriptions/sumo"
	sub, err := client.CreateSubscription(ctxl, "sumo", pubsub.SubscriptionConfig{Topic: top})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sub)

	broker.DefaultBroker = broker.NewBroker(ctx, broker.ProjectID("pro1"), broker.ClientOption(option.WithGRPCConn(conn)))
	return srv
}
