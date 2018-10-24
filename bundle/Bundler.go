package bundle

import (
	"sort"
	"strings"
	"swarm/config"
	"swarm/source"
)

// Bundler is
type Bundler struct {
}

// NewBundler returns a new Bundler
func NewBundler() *Bundler {
	return &Bundler{}
}

// ByFilepath a type to sort files by their names.
type ByFilepath []*source.File

func (nf ByFilepath) Len() int      { return len(nf) }
func (nf ByFilepath) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByFilepath) Less(i, j int) bool {
	nameA := nf[i].Filepath
	nameB := nf[j].Filepath
	return nameA < nameB
}

// Bundle concatenates files in a FileSet into a single file
func (b *Bundler) Bundle(fileset *source.FileSet, runtimeConfig *config.RuntimeConfig) string {
	var sb strings.Builder
	files := fileset.Files()
	sort.Sort(ByFilepath(files))
	for _, file := range files {
		file.EnsureLoaded(runtimeConfig)
		for _, line := range file.BundleBody() {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
