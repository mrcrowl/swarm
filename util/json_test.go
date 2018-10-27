package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonEncodeString(t *testing.T) {
	source := `<abcd>
'efgh'
"ijkl"`

	encoded := JSONEncodeString(source)
	assert.Equal(t, `"<abcd>\n'efgh'\n\"ijkl\""`, encoded)
}
