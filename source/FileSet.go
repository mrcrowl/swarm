package source

import (
	"fmt"
	"sort"
)

// FileSet is
type FileSet struct {
	index        map[string]*File
	links        map[string][]string
	reverseLinks map[string][]string
	workspace    *Workspace
	dirty        bool
}

// NewEmptyFileSet creates an empty FileSet
func NewEmptyFileSet(workspace *Workspace) *FileSet {
	fs := &FileSet{
		index:        make(map[string]*File),
		links:        make(map[string][]string),
		reverseLinks: make(map[string][]string),
		workspace:    workspace,
		dirty:        true,
	}
	return fs
}

// NewFileSet creates a new FileSet initialised with a series of imports and links
func NewFileSet(imports []*Import, links []*DependencyLink, workspace *Workspace) *FileSet {
	fs := NewEmptyFileSet(workspace)
	fs.Ingest(imports, links, false)
	return fs
}

// Ingest extends a FileSet with additional imports
func (fs *FileSet) Ingest(imports []*Import, links []*DependencyLink, replace bool) {
	for _, imp := range imports {
		file, err := fs.workspace.ReadSourceFile(imp)
		if err != nil {
			fmt.Printf("Could not read '%s'\n", imp.Path())
			continue
		}

		if replace {
			fs.Replace(file)
		} else {
			fs.Add(file)
		}
	}

	for _, link := range links {
		fs.AddLink(link)
	}
}

// Dirty gets a flag indicating whether the FileSet needs to be rebundled
func (fs *FileSet) Dirty() bool { return fs.dirty }

// MarkDirty sets a FileSet to dirty
func (fs *FileSet) MarkDirty() { fs.dirty = true }

// ClearDirty clears a FileSet from being dirty
func (fs *FileSet) ClearDirty() { fs.dirty = false }

// Workspace gets the workspace used by this FileSet
func (fs *FileSet) Workspace() *Workspace {
	return fs.workspace
}

// Files returns a list of all Files in the set
func (fs *FileSet) Files() []*File {
	result := make([]*File, 0, len(fs.index))
	for _, id := range fs.calcBundleOrder() {
		if file, found := fs.index[id]; found {
			result = append(result, file)
		} else {
			fmt.Printf("WARN: Files could not find file %s\n", id)
		}
	}
	return result
}

// Get gets a File from the FileSet
func (fs *FileSet) Get(id string) *File /* may be nil */ {
	file := fs.index[id]
	return file
}

// Add adds a File to a FileSet, if it doesn't already contain it
func (fs *FileSet) Add(file *File) bool {
	if fs.Contains(file.ID) {
		return false
	}

	fs.Replace(file)
	return true
}

// Replace overwrites a File in a FileSet
func (fs *FileSet) Replace(file *File) {
	fs.index[file.ID] = file
}

// AddLink adds a DependencyLink between Files in a FileSet
func (fs *FileSet) AddLink(link *DependencyLink) bool {
	if !fs.Contains(link.id) {
		fmt.Printf("ERROR: AddLink() dependent file doesn't exist in the FileSet, ID: %s\n", link.id)
		return false
	}

	for _, dependencyID := range link.dependencyIDs {
		if !fs.Contains(dependencyID) {
			fmt.Printf("ERROR: AddLink() dependency file doesn't exist in the FileSet, ID: %s\n", dependencyID)
			return false
		}
	}

	fs.links[link.id] = link.dependencyIDs
	for _, dependencyID := range link.dependencyIDs {
		if rlinks, found := fs.reverseLinks[dependencyID]; found {
			foundLinkID := false
			for _, rlink := range rlinks {
				if rlink == link.id {
					foundLinkID = true
					break
				}
			}
			if !foundLinkID {
				fs.reverseLinks[dependencyID] = append(rlinks, link.id)
			}
		} else {
			fs.reverseLinks[dependencyID] = []string{link.id}
		}
	}
	return true
}

// contains tests whether a FileSet contains a file
func (fs *FileSet) containsFile(file *File) bool {
	return fs.Contains(file.ID)
}

// Contains tests whether a FileSet contains a file
func (fs *FileSet) Contains(id string) bool {
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
	graph := NewIDGraph(fs.links)
	topoSortedIDs := graph.SortTopologically(fs.sortedFileIDs())
	// if len(leftOverIDs) > 0 {
	// 	graph.analyseLeftoverIDs(leftOverIDs)
	// }

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
