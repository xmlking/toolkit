package broker_test

import (
	"context"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"

	"google.golang.org/api/option"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"

	broker "github.com/xmlking/toolkit/broker/pubsub"
)

func makeMockMessage(newDoneFunc func(string, bool, time.Time)) pubsub.Message {
	message := pubsub.Message{
		ID:              "1",
		Data:            []byte("ABC"),
		Attributes:      map[string]string{"att1": "val1"},
		PublishTime:     time.Now(),
		DeliveryAttempt: nil,
		OrderingKey:     "1",
	}

	//Get a reflectable value of message
	messageValue := reflect.ValueOf(message)

	// The value above is unaddressable. So construct a new and addressable message and set it with the value of the unaddressable
	addressableValue := reflect.New(messageValue.Type()).Elem()
	addressableValue.Set(messageValue)

	//Get message's doneFunc field
	doneFuncField := addressableValue.FieldByName("doneFunc")

	//Get the address of the field
	doneFuncFieldAddress := doneFuncField.UnsafeAddr()

	//Create a pointer based on the address
	doneFuncFieldPointer := unsafe.Pointer(doneFuncFieldAddress)

	//Create a new, exported field element that points to the original
	accessibleDoneFuncField := reflect.NewAt(doneFuncField.Type(), doneFuncFieldPointer).Elem()

	//Set the field with the alternative doneFunc
	accessibleDoneFuncField.Set(reflect.ValueOf(newDoneFunc))

	return addressableValue.Interface().(pubsub.Message)
}

func TestPubsubMessage_Ack(t *testing.T) {
	//Create an alternative done function
	newDoneFunc := func(ackID string, ack bool, receiveTime time.Time) {
		t.Logf("Hi! %s, %t, %v", ackID, ack, receiveTime)
	}
	message := makeMockMessage(newDoneFunc)

	t.Log(message)
	message.Ack()
	t.Log(message)
}

func TestPubsubMessage_Nack(t *testing.T) {
	//Create an alternative done function
	newDoneFunc := func(ackID string, ack bool, receiveTime time.Time) {
		t.Logf("Hi! %s, %t, %v", ackID, ack, receiveTime)
	}
	message := makeMockMessage(newDoneFunc)

	t.Log(message)
	message.Nack()
	t.Log(message)
}

func TestNewBroker(t *testing.T) {
	ctX, cancel := context.WithCancel(context.Background())
	srv := setupFakePubsubAndBroker(ctX, t)
	defer srv.Close()
	defer broker.Shutdown()

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
	if err := broker.Subscribe("sumo", myHandler, broker.WithSubscriptionID("sumo")); err != nil {
		t.Fatal(err)
	}

	// start broker
	if err := broker.Start(); err != nil {
		t.Fatal(err)
	}

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
