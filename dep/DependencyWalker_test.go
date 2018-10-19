package dep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDepedency(t *testing.T) {
	var f = newDependencyQueue()
	f.push("abcd")
	assert.True(t, f.has("abcd"))
	assert.False(t, f.has("efgh"))
}

func TestUniqueDependency(t *testing.T) {
	var f = newDependencyQueue()
	f.push("abcd")
	assert.Equal(t, 1, f.count())
	f.push("abcd")
	assert.Equal(t, 1, f.count())
}

func TestFollowDependencyGraph(t *testing.T) {
	var ws = NewWorkspace("C:\\WF\\LP\\web\\App")
	followDependencyGraph(ws, "app\\src\\ep\\App.js")
}
