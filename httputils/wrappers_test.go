package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseSizeWrapper(t *testing.T) {
	var size int64 = 12
	mw := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(ResponseSizeWrapper(&size, w), req)
		}
	}
	handler := mw(func(w http.ResponseWriter, req *http.Request) {
		t.Helper()
		_, err := w.Write([]byte("foo"))
		require.NoError(t, err)
		_, err = w.Write([]byte("bar"))
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, res.ContentLength, size)
	res, err = http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, res.ContentLength, size)
}

func TestResponseStatusWrapper(t *testing.T) {
	statusCode := 12
	mw := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(ResponseStatusWrapper(&statusCode, w), req)
		}
	}
	handler := mw(func(w http.ResponseWriter, req *http.Request) {
		t.Helper()
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte("foo"))
		require.NoError(t, err)
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, statusCode)
}
