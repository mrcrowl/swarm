package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonEncodeString(t *testing.T) {
	source := `<abcd>
'efgh'
"ijkl"`

	encoded := jsonEncodeString(source)
	assert.Equal(t, `"\u003cabcd\u003e\n'efgh'\n\"ijkl\""`, encoded)
}
