package datasink

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriterSink_Sink(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var buf bytes.Buffer
		sink := &WriterSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), nil),
			Writer:       &buf,
			ErrorHandler: mockErrorHandler(t, nil),
		}
		sink.Sink(exampleData.exStruct)
		assert.Equal(t, "foo", buf.String())
	})

	t.Run("serializer error", func(t *testing.T) {
		var buf bytes.Buffer
		sink := &WriterSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), assert.AnError),
			Writer:       &buf,
			ErrorHandler: mockErrorHandler(t, assert.AnError),
		}
		sink.Sink(exampleData.exStruct)
		assert.Equal(t, "foo", buf.String())
	})

	t.Run("nil writer", func(t *testing.T) {
		sink := &WriterSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), nil),
			Writer:       nil,
			ErrorHandler: mockErrorHandler(t, fmt.Errorf("writer must not be nil")),
		}
		sink.Sink(exampleData.exStruct)
	})

	t.Run("nil serializer", func(t *testing.T) {
		var buf bytes.Buffer
		sink := &WriterSink{
			Serializer:   nil,
			Writer:       &buf,
			ErrorHandler: mockErrorHandler(t, fmt.Errorf("serializer must not be nil")),
		}
		sink.Sink(exampleData.exStruct)
		assert.Empty(t, buf.String())
	})
}
