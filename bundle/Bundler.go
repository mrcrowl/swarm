package bundle

import (
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

// Bundle concatenates files in a FileSet into a single file
func (b *Bundler) Bundle(fileset *source.FileSet, runtimeConfig *config.RuntimeConfig) string {
	var sb strings.Builder
	files := fileset.Files()
	for _, file := range files {
		file.EnsureLoaded(runtimeConfig)
		for _, line := range file.BundleBody() {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
