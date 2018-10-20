package source

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

// ReadSourceFile loads a source file
func (ws *Workspace) ReadSourceFile(imp *Import) (*File, error) {
	exists := false
	absoluteFilePath := ""
	for _, ext := range []string{"", ".js"} {
		absoluteFilePath = filepath.Join(ws.rootPath, (imp.Path() + ext))
		if _, err := os.Stat(absoluteFilePath); os.IsNotExist(err) {
			continue
		}
		exists = true
		break
	}

	if exists {
		return newFile(imp.Path(), absoluteFilePath), nil
	}

	return nil, os.ErrNotExist
}
