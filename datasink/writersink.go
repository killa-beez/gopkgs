package datasink

import (
	"fmt"
	"io"
)

//WriterSink writes data to a writer
type WriterSink struct {
	Serializer   Serializer
	Writer       io.Writer
	ErrorHandler ErrorHandler
}

//Sink writes data to a writer
func (s *WriterSink) Sink(data interface{}) {
	errHandler := s.ErrorHandler
	if errHandler == nil {
		errHandler = func(error) {}
	}
	if s.Serializer == nil {
		errHandler(fmt.Errorf("serializer must not be nil"))
		return
	}
	if s.Writer == nil {
		errHandler(fmt.Errorf("writer must not be nil"))
		return
	}
	err := s.Serializer(data, s.Writer)
	if err != nil {
		errHandler(err)
		return
	}
}
