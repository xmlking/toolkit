package broker

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"cloud.google.com/go/pubsub"
)

// AckHandler implements ack/nack handling.
type AckHandler interface {
	// OnAck processes a message ack.
	OnAck()

	// OnNack processes a message nack.
	OnNack()
}

type makeAckHandler struct {
	t *testing.T
}

func (ah *makeAckHandler) OnAck() {
	ah.t.Logf("OnAck")
}
func (ah *makeAckHandler) OnNack() {
	ah.t.Logf("OnNack")
}

func makeMockMessage(ackHandler AckHandler, t *testing.T) pubsub.Message {
	t.Helper()
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
	//doneFuncField := addressableValue.FieldByName("doneFunc")
	doneFuncField := addressableValue.FieldByName("ackh")

	//Get the address of the field
	doneFuncFieldAddress := doneFuncField.UnsafeAddr()

	//Create a pointer based on the address
	doneFuncFieldPointer := unsafe.Pointer(doneFuncFieldAddress)

	//Create a new, exported field element that points to the original
	accessibleDoneFuncField := reflect.NewAt(doneFuncField.Type(), doneFuncFieldPointer).Elem()

	//Set the field with the alternative doneFunc
	accessibleDoneFuncField.Set(reflect.ValueOf(ackHandler))

	return addressableValue.Interface().(pubsub.Message)
}

func setup() {
	fmt.Println("Setup...")
}

func teardown() {
	fmt.Println("Teardown...")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestPubsubMessage_Ack(t *testing.T) {
	//Create an alternative AckHandler
	ackHandler := &makeAckHandler{t}
	message := makeMockMessage(ackHandler, t)

	t.Log(message)
	message.Ack()
	t.Log(message)
}

func TestPubsubMessage_Nack(t *testing.T) {
	//Create an alternative done function
	ackHandler := &makeAckHandler{t}
	message := makeMockMessage(ackHandler, t)

	t.Log(message)
	message.Nack()
	t.Log(message)
}
