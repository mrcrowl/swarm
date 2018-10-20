package source

import "log"

// FileSet is
type FileSet struct {
	index map[string]*File
	links map[string][]string
}

// NewEmptyFileSet creates an empty FileSet
func NewEmptyFileSet() *FileSet {
	fs := &FileSet{
		index: make(map[string]*File),
		links: make(map[string][]string),
	}
	return fs
}

// NewFileSet creates a new FileSet initialised with a series of imports and links
func NewFileSet(imports []*Import, links []*DependencyLink, workspace *Workspace) *FileSet {
	fs := NewEmptyFileSet()

	for _, imp := range imports {
		file, err := workspace.ReadSourceFile(imp)
		if err != nil {
			log.Printf("Could not read '%s'\n", imp.Path())
			continue
		}

		fs.Add(file)
	}

	for _, link := range links {
		fs.AddLink(link)
	}

	return fs
}

// Files returns a list of all Files in the set
func (fs *FileSet) Files() []*File {
	result := make([]*File, len(fs.index))
	i := 0
	for _, v := range fs.index {
		result[i] = v
		i++
	}
	return result
}

// Add adds a File to a FileSet
func (fs *FileSet) Add(file *File) bool {
	if fs.contains(file.ID) {
		return false
	}

	fs.index[file.ID] = file
	return true
}

// AddLink adds a DependencyLink between Files in a FileSet
func (fs *FileSet) AddLink(link *DependencyLink) bool {
	if !fs.contains(link.id) {
		log.Printf("AddLink: dependent file doesn't exist in the FileSet, ID: %s\n", link.id)
		return false
	}

	for _, dependencyID := range link.dependencyIDs {
		if !fs.contains(dependencyID) {
			log.Printf("AddLink: dependency file doesn't exist in the FileSet, ID: %s\n", dependencyID)
			return false
		}
	}

	fs.links[link.id] = link.dependencyIDs
	return true
}

// contains tests whether a FileSet contains a file
func (fs *FileSet) containsFile(file *File) bool {
	return fs.contains(file.ID)
}

// contains tests whether a FileSet contains a file
func (fs *FileSet) contains(id string) bool {
	if _, ok := fs.index[id]; ok {
		return true
	}
	return false
}

// Count returns the number of files
func (fs *FileSet) Count() int {
	return len(fs.index)
}

// count returns the number of files
func (fs *FileSet) linkCount() int {
	return len(fs.links)
}

func (fs *FileSet) nonEmpty() bool {
	return fs.Count() > 0
}
