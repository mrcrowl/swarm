package source

import (
	"path/filepath"
	"swarm/config"
	"swarm/io"
)

// File represents a single file containing source code
type File struct {
	ID        string
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
	contents, err := io.ReadContents(file.Filepath)
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
func (file *File) SourceMap() *Mapping {
	if file.contents == nil {
		return nil
	}

	sourceMappingURL := file.contents.SourceMappingURL()
	if sourceMappingURL == "" {
		return nil
	}
	absoluteFilepath := filepath.Join(filepath.Dir(file.Filepath), sourceMappingURL)
	return NewMapping(sourceMappingURL, absoluteFilepath)
}

// BundleBody returns a list of lines from the body ready to include in a SystemJSBundle
func (file *File) BundleBody() []string {
	return file.contents.BundleLines()
}
