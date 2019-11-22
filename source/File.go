package source

import (
	"path/filepath"
	"strings"
	"github.com/mrcrowl/swarm/config"
	"github.com/mrcrowl/swarm/util"
)

// File represents a single file containing source code
type File struct {
	ID        string // also happens to be the root-relative url for this file
	Filepath  string
	ext       string
	contents  FileContents
	sourceMap *Mapping
}

// newFile creates a new SourceFile
func newFile(id string, absoluteFilepath string) *File {
	ext := filepath.Ext(absoluteFilepath)
	return &File{
		ID:       id,
		Filepath: absoluteFilepath,
		ext:      ext,
	}
}

func (file *File) rootRelativeURL() string { return file.ID }

// PathRelativeTo returns a path relative to another path
func (file *File) PathRelativeTo(runtimeConfig *config.RuntimeConfig, anotherPath string) string {
	filePath := file.rootRelativeURL()
	basedAnotherPath, err := filepath.Rel(runtimeConfig.BaseHref, anotherPath)
	relativeFilepath, err := filepath.Rel(basedAnotherPath, filePath)
	if err != nil {
		return file.rootRelativeURL()
	}
	return strings.Replace(relativeFilepath, "\\", "/", -1) + ".ts"
}

// Ext gets a file's extension
func (file *File) Ext() string {
	return file.ext
}

// Loaded gets whether a file's contents are loaded
func (file *File) Loaded() bool {
	return file.contents != nil
}

// EnsureLoaded ensures that the Load method has been called for this File instance
func (file *File) EnsureLoaded(runtimeConfig *config.RuntimeConfig) {
	if !file.Loaded() {
		file.LoadContents(runtimeConfig)
	}
}

// LoadContents loads a file's contents from disk and prepares them for bundling
func (file *File) LoadContents(runtimeConfig *config.RuntimeConfig) {
	contents, err := util.ReadContents(file.Filepath)
	if err != nil {
		file.contents = &FailedFileContents{}
		return
	}

	var baseHref string
	if runtimeConfig != nil {
		baseHref = runtimeConfig.BaseHref
	}

	switch file.ext {
	case ".js":
		file.contents, err = ParseJSFileContents(file.ID, contents)
	case ".css":
		file.contents, err = ParseCSSFileContents(file.ID, contents, baseHref)
	default:
		file.contents, err = ParseStringFileContents(file.ID, contents)
	}

	if file.contents == nil {
		panic("ah!")
	}
}

// UnloadContents clears a file's contents
func (file *File) UnloadContents() {
	file.contents = nil
	file.sourceMap = nil
}

// SourceMap gets a Mapping that wraps the sourceMappingURL found within the file's contents.
// This only returns true if the file's contents have been loaded
func (file *File) SourceMap(runtimeConfig *config.RuntimeConfig, entryPointRootRelativePath string) *Mapping {
	if file.sourceMap == nil {
		if file.contents == nil {
			return nil
		}

		sourceMappingURL := file.contents.SourceMappingURL()
		if sourceMappingURL == "" {
			return nil
		}
		relativePath := file.PathRelativeTo(runtimeConfig, entryPointRootRelativePath)
		absoluteFilepath := filepath.Join(filepath.Dir(file.Filepath), sourceMappingURL)
		file.sourceMap = NewMapping(sourceMappingURL, relativePath, absoluteFilepath)
	}
	return file.sourceMap
}

// BundleBody returns a list of lines from the body ready to include in a SystemJSBundle
func (file *File) BundleBody() []string {
	return file.contents.BundleLines()
}

// RawContents provides access to the underlying file contents object
func (file *File) RawContents() FileContents {
	return file.contents
}
