package httputils

import (
	"net/http"
	"sync"
)

//ResponseSizeWrapper returns a ResponseWriter that wraps w and writes the response size to *size
func ResponseSizeWrapper(size *int64, w http.ResponseWriter) http.ResponseWriter {
	wrapper := NewResponseWrapper(w)
	*size = 0
	var mux sync.Mutex
	wrapper.SetWriteFunc(func(p []byte) (int, error) {
		mux.Lock()
		defer mux.Unlock()
		n, err := wrapper.ResponseWriter.Write(p)
		*size += int64(n)
		return n, err
	})
	return wrapper
}

//ResponseStatusWrapper returns a ResponseWriter that wraps w and writes the status code to *statusCode
func ResponseStatusWrapper(statusCode *int, w http.ResponseWriter) http.ResponseWriter {
	var once sync.Once
	wrapper := NewResponseWrapper(w)
	*statusCode = http.StatusOK
	var mux sync.Mutex
	var hasWritten bool
	wrapper.SetWriteFunc(func(p []byte) (int, error) {
		mux.Lock()
		defer mux.Unlock()
		hasWritten = true
		return wrapper.ResponseWriter.Write(p)
	})
	wrapper.SetWriteHeaderFunc(func(code int) {
		mux.Lock()
		defer mux.Unlock()
		if !hasWritten {
			once.Do(func() {
				*statusCode = code
			})
		}
		wrapper.ResponseWriter.WriteHeader(code)
	})
	return wrapper
}
