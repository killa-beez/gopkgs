package datasink

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONSerializer_Serialize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := exampleData.exStruct
		serializer := &JSONSerializer{
			Indent: "  ",
		}
		want := `{
  "Foo": "foo",
  "Bar": "bar"
}
`
		var buf bytes.Buffer
		err := serializer.Serialize(data, &buf)
		assert.NoError(t, err)
		assert.Equal(t, want, buf.String())
	})

	t.Run("nil data", func(t *testing.T) {
		serializer := &JSONSerializer{}
		var buf bytes.Buffer
		err := serializer.Serialize(nil, &buf)
		assert.NoError(t, err)
		assert.Equal(t, "null\n", buf.String())
	})

	t.Run("nil writer", func(t *testing.T) {
		serializer := &JSONSerializer{}
		err := serializer.Serialize(exampleData.exStruct, nil)
		assert.Error(t, err)
		assert.Equal(t, "writer must not be nil", err.Error())
	})
}
