package codes

import (
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
)

// A Code is an unsigned 32-bit error code as defined in the gRPC spec.
type Code uint32

const (
	// OK is returned on success.
	OK Code = 0

	// Unknown is error that do not have enough error information
	Unknown Code = 2

	// TemporaryErrors - such as network problems, server unavailability, etc.
	// retryable immediately with exponential backoff
	TempUnavailable Code = 50

	// SystemErrors - such as misconfiguration, unavailable dependencies, etc.
	// retryable after fixing system config or restoring failed services.
	// SystemErrors fall in between 100 - 149
	SystemTokenExpired     Code = 100
	SystemPathDoesNotExist Code = 101

	// DataErrors - such as invalid/missing input, insufficient/corrupted payload, etc.
	// not-retryable as they are permanent errors. fix the payload and replay.
	// DataErrors fall in between 150 - 199
	DataSchemaNotFound Code = 150
	DataInvalidInput   Code = 151

	_maxCode = 200
)

var strToCode = map[string]Code{
	`"OK"`:      OK,
	`"UNKNOWN"`: Unknown,

	`"TEMP_UNAVAILABLE"`: TempUnavailable,

	`"SYSTEM_TOKEN_EXPIRED"`:       SystemTokenExpired,
	`"SYSTEM_PATH_DOES_NOT_EXIST"`: SystemPathDoesNotExist,

	`"DATA_SCHEMA_NOT_FOUND"`: DataSchemaNotFound,
	`"DATA_INVALID_INPUT"`:    DataInvalidInput,
}

// UnmarshalJSON unmarshals b into the Code.
func (c *Code) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= _maxCode {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*c = Code(ci)
		return nil
	}

	if jc, ok := strToCode[string(b)]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

func (c Code) String() string {
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
		return "DataSchemaNameNotFound"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
