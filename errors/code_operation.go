package errors

import (
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/xmlking/toolkit/errors/categories"
	"github.com/xmlking/toolkit/errors/codes"
)

/*
klError defines a standard application error.
klError is a Wrapper type error (not leaf)
https://middlemost.com/failure-is-your-domain/
Contains:
1. Machine-readable error code.
2. Human-readable message & details
3. Logical operation and nested error
*/

// ErrorOperation is implemented by types that can provide
// Logical Operation that caused the failure
type ErrorOperation interface {
	error
	Operation() string
}

// ErrorHinter is implemented by types that can provide
// Machine-readable error code.
type ErrorCoder interface {
	error
	Code() codes.Code
	Category() categories.Category
}

// Wrapper constructors are for wrapping additional traits to previous error
// WithCode adds code to an existing error.
func WithCode(err error, code codes.Code) error {
	if err == nil {
		return nil
	}
	return &withCode{cause: err, code: code}
}

// WithOperation adds operation to an existing error.
func WithOperation(err error, operation string) error {
	if err == nil {
		return nil
	}
	return &withOperation{cause: err, operation: operation}
}

// WithCodeAndOperation adds code and operation to an existing error.
func WithCodeAndOperation(err error, code codes.Code, operation string) error {
	if err == nil {
		return nil
	}
	return &withOperation{cause: WithCode(err, code), operation: operation}
}

// New constructors are for creating Leaf errors
func New(c codes.Code, operation string, msg string) error {
	return WithCodeAndOperation(errors.New(msg), c, operation)
}
func Newf(c codes.Code, operation string, format string, a ...interface{}) error {
	return New(c, operation, fmt.Sprintf(format, a...))
}

// Helpers

func GetCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if codeErr := ErrorCoder(nil); errors.As(err, &codeErr) {
		return codeErr.Code()
	}
	var timeouter interface {
		error
		Timeout() bool
	}
	var temper interface {
		error
		Temporary() bool
	}
	if (errors.As(err, &temper) && temper.Temporary()) || (errors.As(err, &timeouter) && timeouter.Timeout()) {
		return codes.TempUnavailable
	}

	return codes.Unknown
}

func GetCategory(err error) categories.Category {
	if err == nil {
		return categories.Unknown
	}
	if codeErr := ErrorCoder(nil); errors.As(err, &codeErr) {
		return codeErr.Category()
	}
	return categories.Unknown
}

func GetOperation(err error) string {
	if err == nil {
		return ""
	}
	if opeErr := ErrorOperation(nil); errors.As(err, &opeErr) {
		return opeErr.Operation()
	}

	return "Unknown"
}
