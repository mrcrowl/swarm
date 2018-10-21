package bundle

import (
	"gospm/source"
	"io/ioutil"
	"os"
	"path"
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
func (b *Bundler) Bundle(fileset *source.FileSet) {
	var sb strings.Builder
	for _, file := range fileset.Files() {
		if path.Ext(file.ID) == "" {
			file.EnsureLoaded()
			for _, line := range file.BundleBody() {
				sb.WriteString(line)
				sb.WriteString("\r\n")
			}
		}
	}

	ioutil.WriteFile("c:\\bundle.js", []byte(sb.String()), os.ModePerm)
}
