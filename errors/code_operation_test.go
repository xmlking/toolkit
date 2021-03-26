package errors

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/stretchr/testify/assert"

	"github.com/xmlking/toolkit/errors/categories"
	"github.com/xmlking/toolkit/errors/codes"
)

func TestWithCode(t *testing.T) {
	assertions := assert.New(t)

	origErr := errors.New("world")

	codeErr := WithCode(origErr, codes.SystemPathDoesNotExist)
	assertions.True(errors.Is(codeErr, origErr))
	t.Logf("code err: %+v", codeErr)

	var typedErr ErrorCoder
	if errors.As(codeErr, &typedErr) {
		t.Log(typedErr.Code())
		t.Log(typedErr.Category())
		t.Log(typedErr.Error())
	}
}

func TestWithCodeAndOperation(t *testing.T) {
	assertions := assert.New(t)

	leafErr := errors.New("world")

	coErr := WithCodeAndOperation(leafErr, codes.DataInvalidInput, "test1")
	assertions.True(errors.Is(coErr, leafErr))

	t.Log(coErr)

	var opeErr ErrorOperation
	var codeErr ErrorCoder

	if errors.As(coErr, &opeErr) {
		t.Log(opeErr.Operation())
	}

	if errors.As(coErr, &codeErr) {
		t.Log(codeErr.Code())
	}

	assertions.True(errors.Is(coErr, leafErr))
	assertions.True(errors.Is(coErr, WithCode(leafErr, codes.DataSchemaNotFound)))
}

func TestNew(t *testing.T) {
	assertions := assert.New(t)

	leafErr := errors.New("world")

	coErr := New(codes.DataSchemaNotFound, "test1", "world")
	assertions.True(errors.Is(coErr, leafErr))

	t.Log(coErr)

	var opeErr ErrorOperation
	var codeErr ErrorCoder

	if errors.As(coErr, &opeErr) {
		t.Log(opeErr.Operation())
	}

	if errors.As(coErr, &codeErr) {
		t.Log(codeErr.Code())
	}

	assertions.True(errors.Is(coErr, leafErr))
	assertions.True(errors.Is(coErr, WithCode(leafErr, codes.DataSchemaNotFound)))
}

func TestNewf(t *testing.T) {
	assertions := assert.New(t)

	leafErr := errors.New("hello world")

	coErr := Newf(codes.DataSchemaNotFound, "", "hello %s", "world")
	assertions.True(errors.Is(coErr, leafErr))

	t.Log(coErr)

	var opeErr ErrorOperation
	var codeErr ErrorCoder

	if errors.As(coErr, &opeErr) {
		t.Log(opeErr.Operation())
	}

	if errors.As(coErr, &codeErr) {
		t.Log(codeErr.Code())
	}

	assertions.True(errors.Is(coErr, leafErr))
	assertions.True(errors.Is(coErr, WithCode(leafErr, codes.DataSchemaNotFound)))
}

func TestGetCode(t *testing.T) {
	assertions := assert.New(t)

	coErr := Newf(codes.DataSchemaNotFound, "test", "hello %s", "world")
	assertions.Equal(codes.DataSchemaNotFound, GetCode(coErr))
}

func TestGetCategory(t *testing.T) {
	assertions := assert.New(t)

	coErr := Newf(codes.TempUnavailable, "test", "hello %s", "world")
	assertions.Equal(codes.TempUnavailable, GetCode(coErr))
	assertions.Equal(categories.Temporary, GetCategory(coErr))
}

func TestGetOperation(t *testing.T) {
	assertions := assert.New(t)

	coErr := Newf(codes.DataSchemaNotFound, "test", "hello %s", "world")
	assertions.Equal("test", GetOperation(coErr))

	coErrWithHint := errors.WithHint(coErr, "sumo hint")
	t.Log(errors.FlattenHints(coErrWithHint))

	assertions.Equal("test", GetOperation(coErrWithHint))
	assertions.Equal(codes.DataSchemaNotFound, GetCode(coErrWithHint))
	assertions.Equal(categories.Data, GetCategory(coErrWithHint))
	assertions.Equal([]string{"sumo hint"}, errors.GetAllHints(coErrWithHint))

	errV := fmt.Sprintf("%+v", coErrWithHint)
	assertions.Contains(errV, "operation: test")
	assertions.Contains(errV, "code: DataSchemaNotFound")
	assertions.True(errors.HasType(coErrWithHint, coErr))
	assertions.True(errors.HasType(coErrWithHint, &withCode{}))
}

func TestDecodeError(t *testing.T) {
	assertions := assert.New(t)
	origErr := Newf(codes.DataSchemaNotFound, "test", "hello %s", "world")

	enc := errbase.EncodeError(context.Background(), origErr)
	newErr := errbase.DecodeError(context.Background(), enc)

	t.Logf("end err: %+v", newErr)

	// Ensure that the decorated error can be found.
	// This checks that the wrapper identity
	// is properly preserved across the network.
	assertions.True(errors.Is(newErr, origErr))
}

func ExampleNewf() {
	coErr := Newf(codes.DataSchemaNotFound, "test", "hello %s", "world")

	coErrWithHint := errors.WithHint(coErr, "sumo hint")
	// fmt.Printf("%+v", coErrWithHint)

	fmt.Println(errors.FlattenHints(coErrWithHint))
	fmt.Println(GetOperation(coErrWithHint))
	fmt.Println(GetCode(coErrWithHint))
	fmt.Println(GetCategory(coErrWithHint))
	fmt.Println(errors.GetAllHints(coErrWithHint))

	// Output:
	//
	// sumo hint
	// test
	// DataSchemaNotFound
	// Data
	// [sumo hint]
}
