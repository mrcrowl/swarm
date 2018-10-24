package source

// FileContents is an interface to file contents prepared for SystemJS bundling
type FileContents interface {
	BundleLines() []string
	SourceMappingURL() string
}

// FailedFileContents describes a file that failed to load
type FailedFileContents struct {
}

// BundleLines returns nil
func (ffc *FailedFileContents) BundleLines() []string {
	return nil
}

// SourceMappingURL returns ""
func (ffc *FailedFileContents) SourceMappingURL() string {
	return ""
}
