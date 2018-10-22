package dep

import (
	"swarm/source"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestFollowDependencyGraph(t *testing.T) {
	var ws = source.NewWorkspace("C:\\WF\\LP\\web\\App")
	followDependencyChain(ws, "app\\src\\ep\\App.js", nil)
}

func TestRelativePathToID(t *testing.T) {
	id := relativePathToID("abcd/defg/ghij.js")
	assert.Equal(t, "abcd/defg/ghij", id)
}
