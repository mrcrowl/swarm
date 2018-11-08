package source

import (
	"strings"
)

const systemJSRegisterPrefix = "System.register(["

// ParseRegisterDependencies parses the first line of a SystemJS formatted file and returns the import dependencies
func ParseRegisterDependencies(line string, trimQuotes bool) ([]string, bool) {
	if !strings.HasPrefix(line, systemJSRegisterPrefix) {
		return nil, false // not a register line
	}

	openPos := strings.Index(line, "[")
	closePos := strings.LastIndex(line, "]")
	if openPos < 0 || closePos < 0 || openPos > closePos {
		return nil, false // not a register line
	}

	if closePos == (openPos + 1) {
		return nil, true // no imports
	}

	dependencySlice := line[(openPos + 1):closePos]
	dependencies := strings.Split(dependencySlice, ", ")
	if trimQuotes {
		for i, quotedDependency := range dependencies {
			dependencies[i] = strings.Trim(quotedDependency, "\"")
		}
	}

	return dependencies, true // has imports
}

func skipPreamble(lines []string) ([]string, int) {
	n := len(lines)
	i := 0
	inBlockComment := false
	for ; i < n; i++ {
		line := lines[i]
		if strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "/*") {
			inBlockComment = true
		}
		if inBlockComment {
			if strings.Contains(line, "*/") {
				inBlockComment = false
			}
			continue
		}

		break
	}

	preambleLines := lines[:i]

	return preambleLines, i
}

func chooseBodyLines(lines []string, numPreambleLines int, foundSourceMap bool) []string {
	if foundSourceMap {
		return lines[numPreambleLines : len(lines)-1]
	}
	return lines[numPreambleLines:]
}

const sourceMappingURLPrefix = "//# sourceMappingURL="

// parseSourceMappingURL extracts the sourceMappingURL from a line of text
func parseSourceMappingURL(line string) (string, bool) {
	if strings.HasPrefix(line, sourceMappingURLPrefix) {
		return line[len(sourceMappingURLPrefix):], true
	}
	return "", false
}
