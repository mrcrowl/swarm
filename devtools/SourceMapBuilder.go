package devtools

import (
	"log"
	"strings"
)

// SourceMapBuilder is used for compiling source maps from existing source map files
type SourceMapBuilder struct {
	filename string
	sources  []*sourceMap
}

// NewSourceMapBuilder creates a new sourceMapBuilder
func NewSourceMapBuilder(filename string, capacity int) *SourceMapBuilder {
	return &SourceMapBuilder{
		filename: filename,
		sources:  make([]*sourceMap, 0, capacity),
	}
}

// AddSourceMap adds a source map to be included in the build
func (smb *SourceMapBuilder) AddSourceMap(line int, lineCount int, path string, sourceMapContents string) {
	source := &sourceMap{ // TODO         ^^^^^^^^^^^^^^^^^^^^^^^ remove
		line,
		lineCount,
		path,
		sourceMapContents,
	}
	smb.sources = append(smb.sources, source)
}

func (smb *SourceMapBuilder) String() string {
	var sb strings.Builder
	sb.WriteString(`{"version":3,"file":"`)
	sb.WriteString(smb.filename)
	sb.WriteString(`.js","sources":[`)
	first := true
	for _, source := range smb.sources {
		if first {
			first = false
		} else {
			sb.WriteByte(',')
		}
		sb.WriteString("\"" + source.path + "\"")
	}
	sb.WriteString(`],"mappings":"`)
	sb.WriteString(smb.GenerateMappings())
	sb.WriteString(`"}`)
	return sb.String()
}

// GenerateMappings outputs a string of the compiled sourcemap
func (smb *SourceMapBuilder) GenerateMappings() string {
	fileIndex := 0
	var sb strings.Builder
	for i, source := range smb.sources {
		smap, err := ParseSourceMapJSON(source.contents)
		if err != nil {
			log.Printf("ERROR parsing source map: " + source.path)
			continue
		}
		offset := 1
		if i == 0 {
			offset = 0
		}
		mappings := smap.OffsetMappingsSourceFileIndex(offset)
		sb.WriteString(mappings)
		sb.WriteByte(';')
		fileIndex++
	}
	return sb.String()
}
