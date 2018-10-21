package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDGraph(t *testing.T) {
	links := map[string][]string{
		"a": []string{"b", "c"},
		"b": []string{"c"},
		"c": []string{"d"},
	}
	g := newIDGraph(links)
	topoOrder, leftOvers := g.sortTopologically([]string{"a", "b", "c", "d"})
	assert.True(t, assert.ObjectsAreEqual([]string{"d", "c", "b", "a"}, topoOrder))
	assert.Empty(t, leftOvers)
}
