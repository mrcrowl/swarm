package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const ConfigAppID = "app/src/ConfigApp.js"
const ConfigApp = "c:\\wf\\lp\\web\\App\\app\\src\\ConfigApp.js"

func TestContains(t *testing.T) {
	sut := NewEmptyFileSet()
	assert.False(t, sut.contains(ConfigAppID))
	file := newFile(ConfigAppID, ConfigApp)
	sut.Add(file)
	assert.True(t, sut.contains(ConfigAppID))
	assert.True(t, sut.containsFile(file))
	assert.False(t, sut.contains("aksdjfhzzaksjdfh"))
	assert.Equal(t, 1, sut.count())
	assert.True(t, sut.nonEmpty())
	assert.Equal(t, 0, sut.linkCount())
}

func TestDependency(t *testing.T) {
	sut := NewEmptyFileSet()
	sut.Add(newFile("abcd", "c:\\abcd"))
	sut.Add(newFile("abcd", "c:\\abcd"))
	assert.Equal(t, 1, sut.count())
}

func TestAddDistinct(t *testing.T) {
	sut := NewEmptyFileSet()
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
	assert.Equal(t, 3, sut.count())
	assert.Equal(t, 1, sut.linkCount())
}
