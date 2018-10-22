package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormaliseFilepath(t *testing.T) {
	ws := NewWorkspace("c:\\wf\\lp\\web\\App")
	assert.Equal(t, "c:\\wf\\lp\\web\\App\\", ws.rootPath)
}
func TestNormaliseFilepathWindowsButInputWithUnix(t *testing.T) {
	ws := NewWorkspace("c:/wf/lp/web/App")
	assert.Equal(t, "c:\\wf\\lp\\web\\App\\", ws.rootPath)
}

func TestNormaliseFilepathUnix(t *testing.T) {
	emulateUnix()
	ws := NewWorkspace("/usr/wf/lp/web/App")
	assert.Equal(t, "/usr/wf/lp/web/App/", ws.rootPath)
}

func TestNormaliseFilepathUnixButInputWithWindows(t *testing.T) {
	emulateUnix()
	ws := NewWorkspace("\\usr\\wf\\lp\\web\\App\\")
	assert.Equal(t, "/usr/wf/lp/web/App/", ws.rootPath)
}

func TestToRelativePath(t *testing.T) {
	ws := NewWorkspace("c:\\wf\\lp\\web\\App")
	relative, ok := ws.ToRelativePath("c:\\wf\\lp\\web\\App\\app\\src\\ep\\app.js")
	assert.Equal(t, "app/src/ep/app.js", relative)
	assert.True(t, ok)
}

func TestToRelativePathInvalid(t *testing.T) {
	ws := NewWorkspace("c:\\wf\\lp\\web\\App")
	relative, ok := ws.ToRelativePath("c:\\wf\\home\\topo\\topo.js")
	assert.Equal(t, "", relative)
	assert.False(t, ok)
}
