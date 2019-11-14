package datasink

import "io"

//ErrorHandler handles errors
type ErrorHandler func(err error)

//A Sink does something with your data
type Sink func(data interface{})

//Serializer turns data into []byte and writes it to writer
type Serializer func(data interface{}, writer io.Writer) error
