package source

import (
	"swarm/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonEncodeString(t *testing.T) {
	source := `<abcd>
'efgh'
"ijkl"`

	encoded := jsonEncodeString(source)
	assert.Equal(t, `"<abcd>\n'efgh'\n\"ijkl\""`, encoded)
}

func TestCountLines(t *testing.T) {
	source := "abcd\nefgh"
	count, err := countLines(source)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

func TestCountLinesWindows(t *testing.T) {
	source := "abcd\r\nefgh"
	count, err := countLines(source)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

func TestCountLooooongLines(t *testing.T) {
	source := testutil.ReadTextFile("c:\\wf\\lp\\web\\App\\node_modules\\systemjs\\dist", "system.js")
	count, err := countLines(source)
	assert.Nil(t, err)
	assert.Equal(t, 6, count)
}

func TestStringToLines(t *testing.T) {
	source := "abcd\nefgh"
	lines := stringToLines(source)
	assert.Equal(t, 2, len(lines))
}
