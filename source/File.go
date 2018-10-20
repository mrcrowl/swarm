package source

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
