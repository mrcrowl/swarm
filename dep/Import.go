package dep

import (
	"log"
	"path"
	"strings"
)

// Import is an import in one file, possibly root relative
type Import struct {
	filename         string
	directory        string
	directive        string
	isSelfRelative   bool
	isParentRelative bool
	isRooted         bool
	isSolo           bool
}

func newImport(importPath string) *Import {

	isSelfRelative := strings.HasPrefix(importPath, "./")
	isParentRelative := strings.HasPrefix(importPath, "../")
	isRooted := !isSelfRelative && !isParentRelative

	directory := path.Dir(importPath) // either ., prefixed with ../, or of the form abcd/efgh
	if isSelfRelative && directory != "." {
		directory = "./" + directory
	}

	isSolo := false
	if isRooted && directory == "." {
		directory = ""
		isSolo = true
	}

	directive := ""
	if directivePos := strings.Index(importPath, "#?"); directivePos >= 0 {
		directive = importPath[directivePos:]
		importPath = importPath[:directivePos]
	}

	filename := path.Base(importPath)

	return &Import{filename, directory, directive, isSelfRelative, isParentRelative, isRooted, isSolo}
}

func (imp *Import) path() string {
	if imp.isSolo {
		return imp.filename
	}
	return imp.directory + "/" + imp.filename
}

func (imp *Import) toRootRelativeImport(relativeImport *Import) *Import {
	if imp.isRooted {
		if !relativeImport.isRooted {
			importPathRelativeToRoot := path.Join(imp.directory, relativeImport.path())
			return newImport(importPathRelativeToRoot)
		}
		log.Fatalf("toRootRelativeImportPath called with non-relative import: %s\n", relativeImport.path())
	} else {

	}
	log.Fatalf("toRootRelativeImportPath called on non-root-relative import: %s\n", imp.path())
	return nil
}
