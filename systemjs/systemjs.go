package systemjs

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

	lastLineIndex := len(lines) - 1
	registerLine := lines[0]
	sourceMapLine := lines[lastLineIndex]

	imports, foundRegister := ParseRegisterDependencies(registerLine)
	sourceMappingURL, foundSourceMap := parseSourceMappingURL(sourceMapLine)

	body := chooseBodyLines(lines, foundSourceMap)

	return &FileElements{
		name:             "",
		imports:          imports,
		body:             body,
		sourceMappingURL: sourceMappingURL,
		lineCount:        lastLineIndex + 1,
		isSystemJS:       foundRegister,
	}, nil
}

func chooseBodyLines(lines []string, foundSourceMap bool) []string {
	if foundSourceMap {
		return lines[:len(lines)-1]
	}
	return lines
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
func ParseRegisterDependencies(line string) ([]string, bool) {
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
	for i, quotedDependency := range dependencies {
		dependencies[i] = strings.Trim(quotedDependency, "\"")
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
