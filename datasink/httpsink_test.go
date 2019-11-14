package datasink

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(handler func(rw http.ResponseWriter, req *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestHTTPPostSink_Sink(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		done := make(chan struct{})
		server := testServer(func(rw http.ResponseWriter, req *http.Request) {
			t.Helper()
			assert.Equal(t, "/foo/bar", req.URL.Path)
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			_, err := rw.Write([]byte("OK"))
			require.NoError(t, err)
			close(done)
		})
		defer server.Close()
		sink := &HTTPPostSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), nil),
			URL:          server.URL + "/foo/bar",
			ContentType:  "application/json",
			ErrorHandler: mockErrorHandler(t, nil),
		}
		sink.Sink(exampleData.exStruct)
		<-done
	})

	t.Run("nil serializer", func(t *testing.T) {
		sink := &HTTPPostSink{
			Serializer:   nil,
			ErrorHandler: mockErrorHandler(t, fmt.Errorf("serializer must not be nil")),
		}
		sink.Sink(exampleData.exStruct)
	})

	t.Run("serializer error", func(t *testing.T) {
		sink := &HTTPPostSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), assert.AnError),
			ErrorHandler: mockErrorHandler(t, assert.AnError),
		}
		sink.Sink(exampleData.exStruct)
	})

	t.Run("bad url", func(t *testing.T) {
		sink := &HTTPPostSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), nil),
			ErrorHandler: mockErrorHandler(t, errAny),
			URL:          "http://not.a.real.host",
		}
		sink.Sink(exampleData.exStruct)
	})

	t.Run("404", func(t *testing.T) {
		done := make(chan struct{})
		server := testServer(func(rw http.ResponseWriter, req *http.Request) {
			t.Helper()
			assert.Equal(t, "/foo/bar", req.URL.Path)
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			rw.WriteHeader(http.StatusNotFound)
			_, err := rw.Write([]byte("not found"))
			require.NoError(t, err)
			close(done)
		})
		defer server.Close()
		sink := &HTTPPostSink{
			Serializer:   mockSerializer(t, exampleData.exStruct, []byte("foo"), nil),
			URL:          server.URL + "/foo/bar",
			ContentType:  "application/json",
			ErrorHandler: mockErrorHandler(t, fmt.Errorf("http status code: 404")),
		}
		sink.Sink(exampleData.exStruct)
		<-done
	})
}
