package source

import (
	"gospm/io"
	"path/filepath"
)

// File represents a single file containing source code
type File struct {
	ID       string
	Filepath string
	ext      string
	contents FileContents
}

// newFile creates a new SourceFile
func newFile(id string, absoluteFilepath string) *File {
	ext := filepath.Ext(absoluteFilepath)
	return &File{
		ID:       id,
		Filepath: absoluteFilepath,
		ext:      ext,
	}
}

// Loaded gets whether a file's contents are loaded
func (file *File) Loaded() bool {
	return file.contents != nil
}

// EnsureLoaded ensures that the Load method has been called for this File instance
func (file *File) EnsureLoaded() {
	if !file.Loaded() {
		file.LoadContents()
	}
}

// LoadContents loads a file's contents from disk and prepares them for bundling
func (file *File) LoadContents() {
	contents, err := io.ReadContents(file.Filepath)
	if err != nil {
		file.contents = &FailedFileContents{}
		return
	}

	switch file.ext {
	case ".js":
		file.contents, err = ParseJSFileContents(file.ID, contents)
	case ".css":
		fallthrough
	default:
		file.contents, err = ParseStringFileContents(file.ID, contents)
	}
}

// BundleBody returns a list of lines from the body ready to include in a SystemJSBundle
func (file *File) BundleBody() []string {
	return file.contents.BundleLines()
}
