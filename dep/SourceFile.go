package dep

import (
	"bufio"
	"os"
)

// SourceFile represents a single file containing source code
type SourceFile struct {
	filepath         string
	id               string
	dependentFileIDs map[string]bool
}

// NewSourceFile creates a new SourceFile
func NewSourceFile(absoluteFilepath string) *SourceFile {
	return &SourceFile{
		filepath:         absoluteFilepath,
		dependentFileIDs: make(map[string]bool),
	}
}

// ReadFirstLine reads the first line of a text file as a string
func (file *SourceFile) ReadFirstLine() (string, error) {
	f, err := os.OpenFile(file.filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		return sc.Text(), nil
	}

	return "", nil
}
