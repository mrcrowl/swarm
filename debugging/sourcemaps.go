package debugging

import (
	"encoding/json"
	"errors"
	"strings"
)

// SourceMap represents the JSON structure of a source map in .map file
type SourceMap struct {
	Version    int      `json:"version"`
	File       string   `json:"file"`
	SourceRoot string   `json:"sourceRoot"`
	Sources    []string `json:"sources"`
	Names      []string `json:"names"`
	Mappings   string   `json:"mappings"`
}

type line struct {
	segments [][]int
}

func parseMappings(maps string) []*line {
	lineStrings := strings.Split(maps, ";")
	lines := make([]*line, len(lineStrings))
	for i, lineString := range lineStrings {
		lines[i] = parseLineString(lineString)
	}
	return lines
}

func parseLineString(lineString string) *line {
	if lineString == "" {
		return nil
	}
	segmentStrings := strings.Split(lineString, ",")
	segments := make([][]int, len(segmentStrings))
	for i, segmentString := range segmentStrings {
		segments[i] = Decode(segmentString)
	}
	return &line{segments}
}

func parseSourceMapJSON(sourceMapJSON string) (*SourceMap, error) {
	var sm *SourceMap
	err := json.Unmarshal([]byte(sourceMapJSON), &sm)
	if err != nil {
		return nil, errors.New("Invalid JSON in source map: " + err.Error())
	}
	return sm, nil
}
