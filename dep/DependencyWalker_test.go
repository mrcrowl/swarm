package dep

import (
	"swarm/source"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestFollowDependencyGraph(t *testing.T) {
	var ws = source.NewWorkspace("C:\\WF\\LP\\web\\App")
	followDependencyChain(ws, "app\\src\\ep\\App.js")
}
