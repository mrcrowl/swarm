package bundle

import (
	"gospm/source"
	"strings"
)

// Bundler is
type Bundler struct {
}

// NewBundler returns a new Bundler
func NewBundler() *Bundler {
	return &Bundler{}
}

// Bundle concatenates files in a FileSet into a single file
func (b *Bundler) Bundle(fileset *source.FileSet) *strings.Builder {
	var sb strings.Builder
	for _, file := range fileset.Files() {
		file.EnsureLoaded()
		for _, line := range file.BundleBody() {
			sb.WriteString(line)
			sb.WriteString("\r\n")
		}
	}
	return &sb
}
