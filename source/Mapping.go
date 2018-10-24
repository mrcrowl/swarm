package source

import (
	"log"
	"swarm/io"
)

// Mapping is
type Mapping struct {
	sourceMappingURL string
	filepath         string
	contents         string
}

// NewMapping wraps a sourceMappingURL
func NewMapping(sourceMappingURL string, filepath string) *Mapping {
	return &Mapping{sourceMappingURL, filepath, ""}
}

// EnsureLoaded ensures the files contents are loaded
func (mapping *Mapping) EnsureLoaded() {
	if mapping.contents == "" {
		mapping.LoadContents()
	}
}

// Contents returns the contents of the .map file (if loaded)
func (mapping *Mapping) Contents() string {
	return mapping.contents
}

// LoadContents loads the files contents
func (mapping *Mapping) LoadContents() {
	contents, err := io.ReadContents(mapping.filepath)
	if err != nil {
		log.Printf("Failed to load source map: %s", mapping.filepath)
		return
	}

	mapping.contents = contents
}
