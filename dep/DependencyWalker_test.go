package dep

import (
	"gospm/source"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestFollowDependencyGraph(t *testing.T) {
	var ws = source.NewWorkspace("C:\\WF\\LP\\web\\App")
	followDependencyGraph(ws, "app\\src\\ep\\App.js")
}
