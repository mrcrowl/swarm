package dep

import (
	"os"
	"path/filepath"
)

// Workspace is
type Workspace struct {
	rootPath string
}

// NewWorkspace returns a new workspace for the given rootpath
func NewWorkspace(rootPath string) *Workspace {
	return &Workspace{
		rootPath,
	}
}

// RootPath returns the workspace's root path
func (ws *Workspace) RootPath() string {
	return ws.rootPath
}

func (ws *Workspace) readSourceFile(rootRelativeFilepath string) (*SourceFile, error) {
	exists := false
	absoluteFilePath := ""
	for _, ext := range []string{"", ".js"} {
		absoluteFilePath = filepath.Join(ws.rootPath, (rootRelativeFilepath + ext))
		if _, err := os.Stat(absoluteFilePath); os.IsNotExist(err) {
			continue
		}
		exists = true
		break
	}

	if exists {
		return NewSourceFile(absoluteFilePath), nil
	}

	return nil, os.ErrNotExist
}
