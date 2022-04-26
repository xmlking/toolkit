package codes

import (
	"strconv"

	"google.golang.org/grpc/codes"
)

type Code interface {
	String() string
	Int() uint32
}

// A TCode is an unsigned 32-bit error code as defined in the gRPC spec.
type TCode uint32

const (
	// OK is returned on success.
	OK TCode = 0

	// Unknown is error that do not have enough error information
	Unknown TCode = 2

	// TemporaryErrors - such as network problems, server unavailability, etc.
	// retryable immediately with exponential backoff
	TempUnavailable TCode = 50

	// SystemErrors - such as misconfiguration, unavailable dependencies, etc.
	// retryable after fixing system config or restoring failed services.
	// SystemErrors fall in between 100 - 149
	SystemTokenExpired     TCode = 100
	SystemPathDoesNotExist TCode = 101

	// DataErrors - such as invalid/missing input, insufficient/corrupted payload, etc.
	// not-retryable as they are permanent errors. fix the payload and replay.
	// DataErrors fall in between 150 - 199
	DataSchemaNotFound   TCode = 150
	DataInvalidInput     TCode = 151
	DataResourceNotFound TCode = 152

	_maxCode = 200
)

func (c TCode) String() string {
	// handle gRPC codes
	if c < 50 {
		return codes.Code(c).String()
	}

	switch c {
	case OK:
		return "OK"
	case Unknown:
		return "Unknown"
	case TempUnavailable:
		return "TempUnavailable"
	case SystemTokenExpired:
		return "SystemTokenExpired"
	case SystemPathDoesNotExist:
		return "SystemPathDoesNotExist"
	case DataSchemaNotFound:
		return "DataSchemaNotFound"
	case DataInvalidInput:
		return "DataInvalidInput"
	case DataResourceNotFound:
		return "DataResourceNotFound"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

func (c TCode) Int() uint32 {
	return uint32(c)
}
