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
	IsSelfRelative   bool
	IsParentRelative bool
	IsRooted         bool
	IsSolo           bool
}

var reStripInterpolationTemplate = regexp.MustCompile(`#\{.*?\}`)

// NewImportWithInterpolation creates an Import for a path, but first interpolates any values
func NewImportWithInterpolation(importPath string, interpolationValues map[string]string) *Import {
	if containsInterpolation(importPath) {
		importPath = performInterpolation(importPath, interpolationValues)
	}
	return NewImport(importPath)
}

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

	filename := path.Base(importPath)

	return &Import{filename, directory, isSelfRelative, isParentRelative, isRooted, IsSolo}
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

// containsInterpolation indicates whether a part contains a SystemJS interpolation directive: #{...}
func containsInterpolation(importPath string) bool {
	return strings.Contains(importPath, "#{")
}

var interpRe = regexp.MustCompile("#{[^}]*}")

// performInterpolation interpolates a string with a set of values
func performInterpolation(importPath string, interpolationValues map[string]string) string {
	result := interpRe.ReplaceAllStringFunc(importPath, func(match string) string {
		inner := match[2 : len(match)-1]
		pipePos := strings.Index(inner, "|")
		if pipePos >= 0 {
			key := inner[pipePos+1:]
			if value, ok := interpolationValues[key]; ok {
				return value
			}
		}
		return ""
	})
	return result
}

var interpValuesRe = regexp.MustCompile("(Config\\.\\w+)\\s*=\\s*([^;]+?)\\s*(?:;|\\/\\*)")

// readInterpolationValues
func readInterpolationValues(moduleName string, configJSLines []string) map[string]string {
	values := map[string]string{}
	combinedContents := strings.Join(configJSLines, "\n")
	matches := interpValuesRe.FindAllStringSubmatch(combinedContents, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if is, stringValue := isJSPrimitive(value); is {
			values[key] = stringValue
		} else if is, condition, whenTrue, whenFalse := isJSTernary(value); is {
			values[key] = interpretTernary(values, condition, whenTrue, whenFalse)
		}
	}
	return values
}

// isJSPrimitive
func isJSPrimitive(value string) (bool, string) {
	switch {
	case value == "true":
		return true, "true"
	case value == "false":
		return true, "false"
	case strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\""):
		return true, value[1 : len(value)-1]
	default:
		return false, ""
	}
}

var ternaryRe = regexp.MustCompile("^\\s*(.*?)\\s*\\?\\s*(.*?)\\s*:\\s*(.*)$")

// isJSTernary
func isJSTernary(value string) (is bool, condition string, whenTrue string, whenFalse string) {
	match := ternaryRe.FindStringSubmatch(value)
	if match != nil {
		return true, match[1], match[2], match[3]
	}

	return false, "", "", ""
}

// interpretTernary
func interpretTernary(values map[string]string, condition string, whenTrue string, whenFalse string) string {
	value := values[condition]
	if value == "true" {
		if is, trueValue := isJSPrimitive(whenTrue); is {
			return trueValue
		}
	}
	if is, falseValue := isJSPrimitive(whenFalse); is {
		return falseValue
	}

	return ""
}
