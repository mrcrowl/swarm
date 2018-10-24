package bundle

import (
	"sort"
	"strconv"
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
func (b *Bundler) Bundle(fileset *source.FileSet, runtimeConfig *config.RuntimeConfig, entryPointFilename string) (javascript string, sourcemap string) {
	var jsBuilder strings.Builder
	mapBuilder := newSourceMapBuilder(entryPointFilename)
	files := fileset.Files()
	lineNumber := 0
	sort.Sort(ByFilepath(files))
	for _, file := range files {
		file.EnsureLoaded(runtimeConfig)
		if sourceMap := file.SourceMap(); sourceMap != nil {
			sourceMap.EnsureLoaded()
			mapBuilder.WriteSection(lineNumber, 0, sourceMap.Contents())
		}
		for _, line := range file.BundleBody() {
			jsBuilder.WriteString(line)
			jsBuilder.WriteString("\n")
			lineNumber++
		}
	}
	javascript = jsBuilder.String()
	sourcemap = mapBuilder.String()
	return
}

type sourceMapBuilder struct {
	sb        *strings.Builder
	seenFirst bool
}

func newSourceMapBuilder(filename string) *sourceMapBuilder {
	sb := &strings.Builder{}
	sb.WriteString(`{"version":3,"file":"`)
	sb.WriteString(filename)
	sb.WriteString(`","sections":[`)
	return &sourceMapBuilder{sb, false}
}

func (smb *sourceMapBuilder) String() string {
	smb.sb.WriteString(`]}`)
	return smb.sb.String()
}

func (smb *sourceMapBuilder) WriteSection(line int, column int, sourceMapContents string) {
	if !smb.seenFirst {
		smb.seenFirst = true
	} else {
		smb.sb.WriteString(",")
	}
	smb.sb.WriteString("\n")
	smb.sb.WriteString(`{"offset":{"line":`)
	smb.sb.WriteString(strconv.Itoa(line))
	smb.sb.WriteString(`,"column":`)
	smb.sb.WriteString(strconv.Itoa(column))
	smb.sb.WriteString(`},"map":`)
	smb.sb.WriteString(sourceMapContents) // <-- the actual sourcemap file we're injecting
}
