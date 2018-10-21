package source

import (
	"log"
	"sort"
	"time"
)

// FileSet is
type FileSet struct {
	index        map[string]*File
	links        map[string][]string
	reverseLinks map[string][]string
	workspace    *Workspace
}

// NewEmptyFileSet creates an empty FileSet
func NewEmptyFileSet() *FileSet {
	fs := &FileSet{
		index:        make(map[string]*File),
		links:        make(map[string][]string),
		reverseLinks: make(map[string][]string),
		workspace:    nil,
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

// Workspace gets the workspace used by this Fileset
func (fs *FileSet) Workspace() *Workspace {
	return fs.workspace
}

// Files returns a list of all Files in the set
func (fs *FileSet) Files() []*File {
	start := time.Now()
	defer log.Printf("Files took %s", time.Since(start))

	result := make([]*File, 0, len(fs.index))
	for _, id := range fs.calcBundleOrder() {
		if file, found := fs.index[id]; found {
			result = append(result, file)
		} else {
			log.Printf("WARN: Files could not find file %s\n", id)
		}
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
	for _, dependencyID := range link.dependencyIDs {
		if rlinks, found := fs.reverseLinks[dependencyID]; found {
			fs.reverseLinks[dependencyID] = append(rlinks, link.id)
		} else {
			fs.reverseLinks[dependencyID] = []string{link.id}
		}
	}
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

func (fs *FileSet) calcBundleOrder() []string {
	graph := newIDGraph(fs.links)
	topoSortedIDs, _ := graph.sortTopologically(fs.sortedFileIDs())
	return topoSortedIDs
}

func (fs *FileSet) sortedFileIDs() []string {
	ids := make([]string, len(fs.index))
	i := 0
	for id := range fs.index {
		ids[i] = id
		i++
	}

	sort.StringSlice(ids).Sort()
	return ids
}

func (fs *FileSet) copyLinks() map[string][]string {
	clone := make(map[string][]string)
	for k, v := range fs.links {
		clone[k] = append([]string(nil), v...)
	}
	return clone
}

func (fs *FileSet) indepdentFileIDs() stringStack {
	independentIDs := make([]string, 0, 256)
	for k := range fs.index {
		dependencies := fs.links[k]
		if len(dependencies) == 0 {
			independentIDs = append(independentIDs, k)
		}
	}
	return independentIDs
}
