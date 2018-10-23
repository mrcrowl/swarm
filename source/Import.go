package source

import (
	"log"
	"path"
	"regexp"
	"strings"
)

// Import is an import in one file, possibly root relative
type Import struct {
	Filename         string
	Directory        string
	Directive        string
	IsSelfRelative   bool
	IsParentRelative bool
	IsRooted         bool
	IsSolo           bool
}

var reStripInterpolationTemplate = regexp.MustCompile(`#\{.*?\}`)

// NewImport creates an Import for a path
func NewImport(importPath string) *Import {
	isSelfRelative := strings.HasPrefix(importPath, "./")
	isParentRelative := strings.HasPrefix(importPath, "../")
	isRooted := !isSelfRelative && !isParentRelative

	directory := path.Dir(importPath) // either ., prefixed with ../, or of the form abcd/efgh
	if isSelfRelative && directory != "." {
		directory = "./" + directory
	}

	IsSolo := false
	if isRooted && directory == "." {
		directory = ""
		IsSolo = true
	}

	// Clean interpolation templates of the form "../path/file#{Config|Config.RELEASE}.html".
	importPath = reStripInterpolationTemplate.ReplaceAllString(importPath, "")

	directive := ""
	if directivePos := strings.Index(importPath, "#?"); directivePos >= 0 {
		directive = importPath[directivePos:]
		importPath = importPath[:directivePos]
	}

	filename := path.Base(importPath)

	return &Import{filename, directory, directive, isSelfRelative, isParentRelative, isRooted, IsSolo}
}

// Ext returns the extension of the Import's filename
func (imp *Import) Ext() string {
	return path.Ext(imp.Filename)
}

// Path is the complete path for this import
func (imp *Import) Path() string {
	if imp.IsSolo {
		return imp.Filename
	}
	return imp.Directory + "/" + imp.Filename
}

// ToRootRelativeImport converts a relative import to a root relative import, based on the current import (assuming it is root-relative itself)
func (imp *Import) ToRootRelativeImport(relativeImport *Import) *Import {
	if relativeImport.IsSolo {
		return relativeImport
	}

	if imp.IsRooted {
		if !relativeImport.IsRooted {
			importPathRelativeToRoot := path.Join(imp.Directory, relativeImport.Path())
			return NewImport(importPathRelativeToRoot)
		}
		log.Fatalf("Import.ToRootRelativeImport called with non-relative import: %s\n", relativeImport.Path())
	}
	log.Fatalf("Import.ToRootRelativeImport called on non-root-relative import: %s\n", imp.Path())
	return nil
}
