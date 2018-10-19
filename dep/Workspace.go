package dep

import (
	"os"
	"path"
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

func (ws *Workspace) resolveRelativeDependency(rootRelativeFilepath string, dependencyRelativePath string) string {
	return path.Join(path.Dir(rootRelativeFilepath), dependencyRelativePath)
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
