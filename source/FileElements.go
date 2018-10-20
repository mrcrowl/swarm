package source

// FileElements describes a systemjs file
type FileElements struct {
	name             string
	imports          []string
	body             []string
	sourceMappingURL string
	lineCount        int
	isSystemJS       bool
}

// Parse reads the contents of a SystemJS formatted JS file into its component parts
func Parse(contents string) *FileElements {
	return nil
}
