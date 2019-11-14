package datasink

import (
	"encoding/json"
	"fmt"
	"io"
)

//JSONSerializer serializes data to json
type JSONSerializer struct {
	Indent       string
	IndentPrefix string
}

//Serialize serializes data to json
func (s *JSONSerializer) Serialize(data interface{}, writer io.Writer) error {
	if writer == nil {
		return fmt.Errorf("writer must not be nil")
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent(s.IndentPrefix, s.Indent)
	return encoder.Encode(data)
}
