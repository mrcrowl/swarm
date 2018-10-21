package source

import (
	"strings"
)

// FileElements describes a systemjs file
type FileElements struct {
	preamble         []string
	imports          []string
	body             []string
	sourceMappingURL string
	lineCount        int
	isSystemJS       bool
}

// FailedFileElements is the default placeholder for a file that couldn't be loaded
func FailedFileElements() *FileElements {
	return &FileElements{
		preamble:         nil,
		imports:          nil,
		body:             nil,
		sourceMappingURL: "",
		lineCount:        0,
		isSystemJS:       false,
	}
}

// BundleBody returns a list of lines from the body ready to include in a SystemJSBundle
func (elems *FileElements) BundleBody(name string) []string {
	bodyCopy := []string(nil)
	bodyCopy = append(bodyCopy, elems.preamble...)
	bodyCopy = append(bodyCopy, elems.body...)
	bodyCopy[len(elems.preamble)] = elems.GetRegisterLineForBundle(name)
	return bodyCopy
}

// GetRegisterLineForBundle outputs the System.register line with a name
func (elems *FileElements) GetRegisterLineForBundle(name string) string {
	imports := strings.Join(elems.imports, ", ")
	return "System.register(\"" + name + ".js\", [" + imports + "], function (exports_1, context_1) {"
}
