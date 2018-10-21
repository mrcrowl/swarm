package source

import (
	"bufio"
	"strings"
)

// ParseSystemJSFormattedFile parses the lines of a SystemJS formatted file into the Elements
func ParseSystemJSFormattedFile(fileContents string) (*FileElements, error) {
	var lines []string
	var err error

	if lines, err = stringToLines(fileContents); err != nil {
		return nil, err
	}

	numLines := len(lines)
	var imports []string
	var foundRegister = false
	var sourceMappingURL = ""
	var foundSourceMap = false
	var body []string
	var preamble []string
	var numPreambleLines int
	if numLines > 0 {
		preamble, numPreambleLines = skipPreamble(lines)
		registerLine := lines[numPreambleLines]
		imports, foundRegister = ParseRegisterDependencies(registerLine, false)

		sourceMapLine := lines[numLines-1]
		sourceMappingURL, foundSourceMap = parseSourceMappingURL(sourceMapLine)

		body = chooseBodyLines(lines, numPreambleLines, foundSourceMap)
	}

	return &FileElements{
		preamble:         preamble,
		imports:          imports,
		body:             body,
		sourceMappingURL: sourceMappingURL,
		lineCount:        numLines,
		isSystemJS:       foundRegister,
	}, nil
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
			continue
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

func stringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}
