package datasink

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exampleStruct struct{ Foo, Bar string }

var exampleData = struct {
	exStruct exampleStruct
}{
	exStruct: exampleStruct{
		Foo: "foo",
		Bar: "bar",
	},
}

var errAny = fmt.Errorf("any error")

func mockErrorHandler(t *testing.T, wantErr error) ErrorHandler {
	t.Helper()
	return func(err error) {
		t.Helper()
		switch wantErr {
		case nil:
			assert.NoError(t, err)
		case errAny:
			assert.Error(t, err)
		default:
			assert.EqualError(t, err, wantErr.Error())
		}
	}
}

func mockSerializer(t *testing.T, expectsData interface{}, writes []byte, returnsError error) Serializer {
	t.Helper()
	return func(data interface{}, writer io.Writer) error {
		t.Helper()
		if expectsData == nil {
			assert.Nil(t, data)
		} else {
			assert.Equal(t, expectsData, data)
		}
		_, err := writer.Write(writes)
		require.NoError(t, err)
		return returnsError
	}
}
