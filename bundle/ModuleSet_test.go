package bundle

import (
	"swarm/config"
	"swarm/source"
	"testing"

	"github.com/stretchr/testify/assert"
)

const buildDescrSampleJSON = `{
	"modules": [
		{
			"name": "abcd/efgh",
			"include": [
				"common/util",
				"common/dict",
				"common/strings"
			]
		},
		{
			"name": "wxyz/zzzz",
			"exclude": [
				"abcd/efgh"
			]
		},
		{
			"name": "stuv/vvvv",
			"exclude": [
				"abcd/efgh",
				"wxyz/zzzz"
			]
		}
	],
    "base": "app/src/"
}`

// duplicated in other places
func createWorkspace() *source.Workspace {
	ws := source.NewWorkspace("c:\\wf\\lp\\web\\App")
	return ws
}

func TestCreateModuleSet(t *testing.T) {
	descr, err := config.LoadBuildDescriptionString(buildDescrSampleJSON)
	assert.Nil(t, err)

	assert.Len(t, descr.Modules, 3)
	set := CreateModuleSet(createWorkspace(), descr.NormaliseModules("c:\\wf\\lp\\web\\App"))
	assert.True(t, assert.ObjectsAreEqual([]string{"abcd/efgh", "wxyz/zzzz", "stuv/vvvv"}, set.names()), "Module order doesn't match")
}

func TestCreateModuleSetFromFile(t *testing.T) {
	descr, err := config.LoadBuildDescriptionFile("c:\\wf\\lp\\web\\App\\build\\systemjs_build_controlpanel.json")
	assert.Nil(t, err)

	assert.True(t, len(descr.Modules) > 10)
	set := CreateModuleSet(createWorkspace(), descr.NormaliseModules("c:\\wf\\lp\\web\\App"))
	assert.Equal(t, "controlPanel/ControlPanel", set.names()[0], "controlPanel/ControlPanel should be the first module")
}
