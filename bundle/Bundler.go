package bundle

import (
	"path"
	"sort"
	"strings"
	"swarm/config"
	"swarm/devtools"
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
func (b *Bundler) Bundle(fileset *source.FileSet, runtimeConfig *config.RuntimeConfig, entryPointPath string) (javascript string, sourcemap string) {
	var jsBuilder strings.Builder
	entryPointFilename := path.Base(entryPointPath)
	mapBuilder := devtools.NewSourceMapBuilder(entryPointFilename, fileset.Count())
	files := fileset.Files()
	lineIndex := 0
	sort.Sort(ByFilepath(files))
	for _, file := range files {
		file.EnsureLoaded(runtimeConfig)
		startingLineIndex := 0
		for _, line := range file.BundleBody() {
			jsBuilder.WriteString(line)
			jsBuilder.WriteString("\n")
			lineIndex++
		}
		if sourceMap := file.SourceMap(runtimeConfig, entryPointPath); sourceMap != nil {
			lineCount := lineIndex - startingLineIndex
			sourceMap.EnsureLoaded()
			mapBuilder.AddSourceMap(lineIndex, lineCount, sourceMap.RelativePath(), sourceMap.Contents())
		}
	}
	javascript = jsBuilder.String()
	sourcemap = mapBuilder.String()
	return
}

// func newSourceMapBuilder(filename string) *sourceMapBuilder {
// 	sb := &strings.Builder{}
// 	sb.WriteString(`{"version":3,"file":"`)
// 	sb.WriteString(filename)
// 	sb.WriteString(`","sections":[`)
// 	return &sourceMapBuilder{sb, false}
// }

// func (smb *sourceMapBuilder) String() string {
// 	smb.sb.WriteString(`]}`)
// 	return smb.sb.String()
// }

// func (smb *sourceMapBuilder) WriteSection(line int, column int, sourceMapContents string) {
// 	if !smb.seenFirst {
// 		smb.seenFirst = true
// 	} else {
// 		smb.sb.WriteString(",")
// 	}
// 	smb.sb.WriteString("\n")
// 	smb.sb.WriteString(`{"offset":{"line":`)
// 	smb.sb.WriteString(strconv.Itoa(line))
// 	smb.sb.WriteString(`,"column":`)
// 	smb.sb.WriteString(strconv.Itoa(column))
// 	smb.sb.WriteString(`},"map":`)
// 	smb.sb.WriteString(sourceMapContents) // <-- the actual sourcemap file we're injecting
// }
