package source

import (
	"gospm/io"
)

// File represents a single file containing source code
type File struct {
	ID       string
	Filepath string
	elements *FileElements
}

// newFile creates a new SourceFile
func newFile(id string, absoluteFilepath string) *File {
	return &File{
		ID:       id,
		Filepath: absoluteFilepath,
	}
}

// EnsureLoaded ensures that the Load method has been called for this File instance
func (file *File) EnsureLoaded() {
	if file.elements == nil {
		file.Load()
	}
}

// Load loads a file from disk and parses the contents
func (file *File) Load() {
	contents, err := io.ReadContents(file.Filepath)
	if err != nil {
		file.elements = FailedFileElements()
	}
	file.elements, err = ParseSystemJSFormattedFile(contents)
}

// Body returns a list of lines representing the body of the file
func (file *File) Body() []string {
	return file.elements.body
}
