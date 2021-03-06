package dep

import (
	"github.com/mrcrowl/swarm/source"
	"github.com/mrcrowl/swarm/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFollowDependencyGraph(t *testing.T) {
	ws := source.NewWorkspace("C:\\WF\\LP\\web\\App")
	followDependencyChain(ws, "app\\src\\ep\\App.js", nil, map[string]string{})
}

const jsFileWithCommentsBeforeSystemRegister = `// the dependency above is required by evaluateVariables() method
// it attaches methods to the window.ep.functions namespace
System.register(["./QuestionStateContext", "./SupportFunctions", "./VariableDependencySorter"], function (exports_1, context_1) {
...`

func TestReadDependencies(t *testing.T) {
	temppath := testutil.CreateTempDir()
	defer testutil.RemoveTempDir(temppath)
	testutil.WriteTextFile(temppath, "VariableEvaluator.js", jsFileWithCommentsBeforeSystemRegister)

	ws := source.NewWorkspace(temppath)
	imp := source.NewImport("./VariableEvaluator.js")
	file, err := ws.ReadSourceFile(imp)
	assert.Nil(t, err)
	dependencies := readDependencies(file, map[string]string{})
	assert.Len(t, dependencies, 3)
}
