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
	g := NewIDGraph(links)
	topoOrder := g.SortTopologically([]string{"a", "b", "c", "d"})
	assert.True(t, assert.ObjectsAreEqual([]string{"d", "c", "b", "a"}, topoOrder))
}

func TestStringStack(t *testing.T) {
	ss := newStringStack([]string{"a", "b", "c"})
	c := ss.pop()
	assert.Equal(t, "c", c)
	b := ss.pop()
	assert.Equal(t, "b", b)
}
