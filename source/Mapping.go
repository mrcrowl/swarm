package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"swarm/util"
)

// Mapping is
type Mapping struct {
	sourceMappingURL string
	relativePath     string
	filepath         string
	config           *MapConfig
	playback         *MapPlayback
}

// Playback is
func (mapping *Mapping) Playback() *MapPlayback {
	return mapping.playback
}

// CachePlayback stores a playback to avoid it being recalculated
func (mapping *Mapping) CachePlayback(playback *MapPlayback) {
	mapping.playback = playback
}

// Mappings returns the string of source mappings
func (mapping *Mapping) Mappings() string {
	if mapping.config == nil {
		fmt.Println("ERROR: unexpected nil in Mapping.Mappings()")
		return ""
	}
	return mapping.config.Mappings
}

// MapPlayback is a cache of the line count and segment delta
type MapPlayback struct {
	LineCount    int
	SegmentDelta Segment
}

// Segment is a mapping between a source file, line and column --> a generated column
type Segment struct {
	GeneratedColumn int
	SourceFile      int
	SourceLine      int
	SourceColumn    int
}

// AdjustForSource prepares a segment for inversion of another segment
func (seg *Segment) AdjustForSource() Segment {
	return Segment{0, seg.SourceFile, -seg.SourceLine, -seg.SourceColumn}
}

// Add adds the values of two segments
func (seg *Segment) Add(other Segment) Segment {
	return Segment{
		seg.GeneratedColumn + other.GeneratedColumn,
		seg.SourceFile + other.SourceFile,
		seg.SourceLine + other.SourceLine,
		seg.SourceColumn + other.SourceColumn,
	}
}

// MapConfig represents the JSON structure of a source map in .map file
type MapConfig struct {
	Version    int      `json:"version"`
	File       string   `json:"file"`
	SourceRoot string   `json:"sourceRoot"`
	Sources    []string `json:"sources"`
	Names      []string `json:"names"`
	Mappings   string   `json:"mappings"`
}

// ParseSourceMapConfig parses a source map from a json string
func ParseSourceMapConfig(sourceMapJSON string) (*MapConfig, error) {
	var sm *MapConfig
	err := json.Unmarshal([]byte(sourceMapJSON), &sm)
	if err != nil {
		return nil, errors.New("Invalid JSON in source map: " + err.Error())
	}
	return sm, nil
}

// NewMapping wraps a sourceMappingURL
func NewMapping(sourceMappingURL string, relativePath string, filepath string) *Mapping {
	return &Mapping{sourceMappingURL, relativePath, filepath, nil, nil}
}

// NewMappingForTesting is ONLY intended for testing purposes
func NewMappingForTesting(config *MapConfig) *Mapping {
	return &Mapping{config: config}
}

// RelativePath returns the path relative to the entry point
func (mapping *Mapping) RelativePath() string {
	return mapping.relativePath
}

// EnsureLoaded ensures the files contents are loaded
func (mapping *Mapping) EnsureLoaded() {
	if mapping.config == nil {
		mapping.LoadConfig()
	}
}

// NOT REQUIRED BECAUSE FILE ABANDONS THE MAPPING COMPLETELY
// // Unload removes the cached config & playback
// func (mapping *Mapping) Unload() {
// 	mapping.config = nil
// 	mapping.playback = nil
// }

// LoadConfig loads the files contents
func (mapping *Mapping) LoadConfig() {
	contents, err := util.ReadContents(mapping.filepath)
	if err != nil {
		log.Printf("Failed to load source map: %s", mapping.filepath)
		return
	}

	smapConfig, err := ParseSourceMapConfig(contents)
	if err != nil {
		log.Printf("ERROR parsing source map: " + mapping.relativePath)
		return
	}
	mapping.config = smapConfig
}
