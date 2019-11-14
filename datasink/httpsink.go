package datasink

import (
	"bytes"
	"fmt"
	"net/http"
)

//HTTPPostSink posts data to a url
type HTTPPostSink struct {
	Serializer   Serializer
	URL          string
	ContentType  string
	ErrorHandler ErrorHandler
	HTTPClient   *http.Client
}

//Sink posts data to a url
func (s *HTTPPostSink) Sink(data interface{}) {
	if s.HTTPClient == nil {
		s.HTTPClient = http.DefaultClient
	}
	errHandler := s.ErrorHandler
	if errHandler == nil {
		errHandler = func(error) {}
	}
	if s.Serializer == nil {
		errHandler(fmt.Errorf("serializer must not be nil"))
		return
	}
	var buf bytes.Buffer
	err := s.Serializer(data, &buf)
	if err != nil {
		errHandler(err)
		return
	}
	res, err := s.HTTPClient.Post(s.URL, s.ContentType, &buf)
	if err != nil {
		errHandler(err)
		return
	}
	err = res.Body.Close()
	if err != nil {
		errHandler(err)
		return
	}
	if res.StatusCode >= 400 {
		errHandler(fmt.Errorf("http status code: %d", res.StatusCode))
		return
	}
}
