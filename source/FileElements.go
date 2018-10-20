package source

// FileElements describes a systemjs file
type FileElements struct {
	imports          []string
	body             []string
	sourceMappingURL string
	lineCount        int
	isSystemJS       bool
}

// FailedFileElements is the default placeholder for a file that couldn't be loaded
func FailedFileElements() *FileElements {
	return &FileElements{
		imports:          nil,
		body:             nil,
		sourceMappingURL: "",
		lineCount:        0,
		isSystemJS:       false,
	}
}
