package errors

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/cockroachdb/errors/errorspb"
	"github.com/gogo/protobuf/proto"
)

// withOperation is implemented by types that can provide
// Logical Operation that caused the failure
type withOperation struct {
	cause     error
	operation string
}

var _ error = (*withOperation)(nil)
var _ ErrorOperation = (*withOperation)(nil)
var _ errbase.SafeDetailer = (*withOperation)(nil)
var _ fmt.Formatter = (*withOperation)(nil)
var _ errbase.Formatter = (*withOperation)(nil)

func (w *withOperation) Operation() string     { return w.operation }
func (w *withOperation) Error() string         { return w.cause.Error() }
func (w *withOperation) Cause() error          { return w.cause }
func (w *withOperation) Unwrap() error         { return w.cause }
func (w *withOperation) SafeDetails() []string { return []string{w.operation} }

func (w *withOperation) Format(s fmt.State, verb rune) { errbase.FormatError(w, s, verb) }
func (w *withOperation) FormatError(p errors.Printer) (next error) {
	if p.Detail() {
		p.Printf("operation: %s", w.operation)
	}
	return w.cause
}

func encodeWithOperation(_ context.Context, err error) (string, []string, proto.Message) {
	w := err.(*withOperation)
	return "", nil, &errorspb.StringPayload{Msg: w.operation}
}

func decodeWithOperation(
	_ context.Context, cause error, _ string, _ []string, payload proto.Message,
) error {
	m, ok := payload.(*errorspb.StringPayload)
	if !ok {
		// If this ever happens, this means some version of the library
		// (presumably future) changed the payload type, and we're
		// receiving this here. In this case, give up and let
		// DecodeError use the opaque type.
		return nil
	}
	return &withOperation{cause: cause, operation: m.Msg}
}

func init() {
	errbase.RegisterWrapperEncoder(errbase.GetTypeKey((*withOperation)(nil)), encodeWithOperation)
	errbase.RegisterWrapperDecoder(errbase.GetTypeKey((*withOperation)(nil)), decodeWithOperation)
}
