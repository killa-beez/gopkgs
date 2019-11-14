package httputils

import (
	"net/http"
)

//ResponseWrapper wraps a http.ResponseWriter
type ResponseWrapper struct {
	ResponseWriter  http.ResponseWriter
	writeHeaderFunc func(statusCode int)
	headerFunc      func() http.Header
	writeFunc       func(p []byte) (int, error)
}

//NewResponseWrapper returns a new *ResponseWrapper
func NewResponseWrapper(responseWriter http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{
		ResponseWriter: responseWriter,
	}
}

//SetWriteFunc sets the function to run when r.Write([]byte]) is called
func (r *ResponseWrapper) SetWriteFunc(fn func(p []byte) (int, error)) {
	r.writeFunc = fn
}

//SetHeaderFunc sets the function to run when r.Header() is called
func (r *ResponseWrapper) SetHeaderFunc(fn func() http.Header) {
	r.headerFunc = fn
}

//SetWriteHeaderFunc sets the function to run when r.WriteHeader(int) is called
func (r *ResponseWrapper) SetWriteHeaderFunc(fn func(statusCode int)) {
	r.writeHeaderFunc = fn
}

//Header implements http.ResponseWriter
func (r *ResponseWrapper) Header() http.Header {
	if r.headerFunc != nil {
		return r.headerFunc()
	}
	return r.ResponseWriter.Header()
}

//Write implements http.ResponseWriter
func (r *ResponseWrapper) Write(p []byte) (int, error) {
	if r.writeFunc != nil {
		return r.writeFunc(p)
	}
	return r.ResponseWriter.Write(p)
}

//WriteHeader implements http.ResponseWriter
func (r *ResponseWrapper) WriteHeader(statusCode int) {
	if r.writeHeaderFunc != nil {
		r.writeHeaderFunc(statusCode)
	}
	r.ResponseWriter.WriteHeader(statusCode)
}
