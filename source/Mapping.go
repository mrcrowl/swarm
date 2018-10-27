package source

import (
	"log"
	"swarm/util"
)

// Mapping is
type Mapping struct {
	sourceMappingURL string
	relativePath     string
	filepath         string
	contents         string
}

// NewMapping wraps a sourceMappingURL
func NewMapping(sourceMappingURL string, relativePath string, filepath string) *Mapping {
	return &Mapping{sourceMappingURL, relativePath, filepath, ""}
}

// EnsureLoaded ensures the files contents are loaded
func (mapping *Mapping) EnsureLoaded() {
	if mapping.contents == "" {
		mapping.LoadContents()
	}
}

// RelativePath returns the path relative to the entry point
func (mapping *Mapping) RelativePath() string {
	return mapping.relativePath
}

// Contents returns the contents of the .map file (if loaded)
func (mapping *Mapping) Contents() string {
	return mapping.contents
}

// LoadContents loads the files contents
func (mapping *Mapping) LoadContents() {
	contents, err := util.ReadContents(mapping.filepath)
	if err != nil {
		log.Printf("Failed to load source map: %s", mapping.filepath)
		return
	}

	mapping.contents = contents
}
