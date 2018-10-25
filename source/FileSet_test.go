package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const ConfigAppID = "app/src/ConfigApp.js"
const ConfigApp = "c:\\wf\\lp\\web\\App\\app\\src\\ConfigApp.js"

// duplicated in other places
func createWorkspace() *Workspace {
	ws := NewWorkspace("c:\\wf\\lp\\web\\App")
	return ws
}

func TestContains(t *testing.T) {
	sut := NewEmptyFileSet(createWorkspace())
	assert.False(t, sut.Contains(ConfigAppID))
	file := newFile(ConfigAppID, ConfigApp)
	sut.Add(file)
	assert.True(t, sut.Contains(ConfigAppID))
	assert.True(t, sut.containsFile(file))
	assert.False(t, sut.Contains("aksdjfhzzaksjdfh"))
	assert.Equal(t, 1, sut.Count())
	assert.True(t, sut.nonEmpty())
	assert.Equal(t, 0, sut.linkCount())
}

func TestDependency(t *testing.T) {
	sut := NewEmptyFileSet(createWorkspace())
	sut.Add(newFile("abcd", "c:\\abcd"))
	sut.Add(newFile("abcd", "c:\\abcd"))
	assert.Equal(t, 1, sut.Count())
}

func TestAddDistinct(t *testing.T) {
	sut := NewEmptyFileSet(createWorkspace())
	sut.Add(newFile("abcd", "c:\\abcd"))
	sut.Add(newFile("efgh", "c:\\efgh"))
	sut.Add(newFile("ijkl", "c:\\ijkl"))

	var success bool
	success = sut.AddLink(NewDependencyLink("abcd", []string{"efgh", "ijkl"}))
	assert.True(t, success)

	sut.AddLink(NewDependencyLink("efgh", []string{"xyzw"}))
	assert.Equal(t, 1, sut.linkCount())

	sut.AddLink(NewDependencyLink("xyzw", []string{"efgh"}))
	assert.Equal(t, 1, sut.linkCount())
}

func TestNewBuilder(t *testing.T) {
	imports := []*Import{
		NewImport("Config"),
		NewImport("app/index.html"),
		NewImport("app/src/ConfigApp"),
	}
	links := []*DependencyLink{
		NewDependencyLink("app/src/ConfigApp", []string{"Config"}),
	}
	ws := NewWorkspace("C:\\WF\\LP\\web\\App")
	sut := NewFileSet(imports, links, ws)
	assert.Equal(t, 3, sut.Count())
	assert.Equal(t, 1, sut.linkCount())
}

// Topologically-sorted builds weren't required for SystemJS after all.
// ==================================================================
// func TestCalcBundleOrder(t *testing.T) {
// 	imports := []*Import{
// 		NewImport("Config"),
// 		NewImport("app/index.html"),
// 		NewImport("app/src/ConfigApp"),
// 		NewImport("app/src/ep/app.css"),
// 		NewImport("app/src/ep/app"),
// 		NewImport("app/src/ep/AppController"),
// 		NewImport("common/Common"),
// 		NewImport("common/time/TimeBoxManager"),
// 	}
// 	links := []*DependencyLink{
// 		NewDependencyLink("app/src/ConfigApp", []string{"Config", "app/index.html", "app/src/ep/app.css", "common/Common", "app/src/ep/app"}),
// 		NewDependencyLink("app/src/ep/app", []string{"app/src/ep/AppController", "common/Common"}),
// 		NewDependencyLink("common/Common", []string{"common/time/TimeBoxManager"}),
// 	}
// 	ws := NewWorkspace("C:\\WF\\LP\\web\\App")
// 	sut := NewFileSet(imports, links, ws)
// 	assert.Equal(t, sut.Count(), 8)
// 	assert.Equal(t, sut.linkCount(), 3)
// 	order := sut.calcBundleOrder()

// 	common := indexOf("common/Common", order)
// 	timeboxmanager := indexOf("common/time/TimeBoxManager", order)
// 	config := indexOf("Config", order)
// 	indexhtml := indexOf("app/index.html", order)
// 	appcss := indexOf("app/src/ep/app.css", order)
// 	appcontroller := indexOf("app/src/ep/AppController", order)
// 	app := indexOf("app/src/ep/app", order)
// 	configapp := indexOf("app/src/ConfigApp", order)

// 	assert.True(t, configapp > config)
// 	assert.True(t, configapp > indexhtml)
// 	assert.True(t, configapp > appcss)
// 	assert.True(t, configapp > common)
// 	assert.True(t, configapp > app)
// 	assert.True(t, app > appcontroller)
// 	assert.True(t, app > common)
// 	assert.True(t, common > timeboxmanager)
// 	assert.True(t, configapp > timeboxmanager)
// 	assert.True(t, app > timeboxmanager)
// }

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
