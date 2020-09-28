package errors

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors/errbase"
	"github.com/cockroachdb/errors/extgrpc"
	"github.com/cockroachdb/redact"
	"github.com/gogo/protobuf/proto"
	gcodes "google.golang.org/grpc/codes"

	"github.com/xmlking/toolkit/errors/categories"
	"github.com/xmlking/toolkit/errors/codes"
)

type withCode struct {
	cause error
	code  codes.Code
}

var _ error = (*withCode)(nil)
var _ ErrorCoder = (*withCode)(nil)
var _ fmt.Formatter = (*withCode)(nil)
var _ errbase.SafeFormatter = (*withCode)(nil)
var _ errbase.SafeDetailer = (*withCode)(nil)

// withCode is a ErrorCoder.
func (w *withCode) Code() codes.Code {
	if w == nil || w.cause == nil {
		return codes.OK
	}
	return w.code
}

// Category returns the error Category. Read-only
func (w *withCode) Category() categories.Category {
	switch code := w.Code(); {
	case code < 50:
		return grpcCategory(code)
	case code >= 50 && code < 100:
		return categories.Temporary
	case code >= 100 && code < 150:
		return categories.System
	case code >= 150 && code < 200:
		return categories.Data
	default:
		return categories.Unknown
	}
}
func grpcCategory(code codes.Code) categories.Category {
	switch gcodes.Code(code) {
	case gcodes.Canceled, gcodes.DeadlineExceeded, gcodes.ResourceExhausted, gcodes.Aborted, gcodes.Unavailable, gcodes.DataLoss:
		return categories.Temporary
	case gcodes.Unknown, gcodes.PermissionDenied, gcodes.FailedPrecondition, gcodes.Internal, gcodes.Unauthenticated:
		return categories.System
	case gcodes.InvalidArgument, gcodes.NotFound, gcodes.AlreadyExists, gcodes.OutOfRange, gcodes.Unimplemented:
		return categories.Data
	default:
		return categories.Unknown
	}
}

// withCode is an error.
func (w *withCode) Error() string { return w.cause.Error() }

// withCode is also a wrapper.
func (w *withCode) Cause() error          { return w.cause }
func (w *withCode) Unwrap() error         { return w.cause }
func (w *withCode) SafeDetails() []string { return []string{w.code.String()} }

// it knows how to format itself.
func (w *withCode) Format(s fmt.State, verb rune) { errbase.FormatError(w, s, verb) }

// SafeFormatter implements errors.SafeFormatter.
// Note: see the documentation of errbase.SafeFormatter for details
// on how to implement this. In particular beware of not emitting
// unsafe strings.
func (w *withCode) SafeFormatError(p errbase.Printer) error {
	if p.Detail() {
		p.Printf("code: %s", redact.Safe(w.code))
	}
	return w.cause
}

// it's an encodable error.
func encodeWithCode(_ context.Context, err error) (string, []string, proto.Message) {
	w := err.(*withCode)
	details := []string{fmt.Sprintf("gRPC %d", w.code)}
	payload := &extgrpc.EncodedGrpcCode{Code: uint32(w.code)}
	return "", details, payload
}

// it's a decodable error.
func decodeWithCode(
	_ context.Context, cause error, _ string, _ []string, payload proto.Message,
) error {
	wp := payload.(*extgrpc.EncodedGrpcCode)
	return &withCode{cause: cause, code: codes.Code(wp.Code)}
}

func init() {
	errbase.RegisterWrapperEncoder(errbase.GetTypeKey((*withCode)(nil)), encodeWithCode)
	errbase.RegisterWrapperDecoder(errbase.GetTypeKey((*withCode)(nil)), decodeWithCode)
}
