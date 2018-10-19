package dep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFirstLine(t *testing.T) {
	f := NewSourceFile("C:\\WF\\LP\\web\\App\\app\\src\\ep\\AppController.js")
	line, err := f.ReadFirstLine()
	assert.Nil(t, err)
	print(line)
}
