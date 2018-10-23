package io

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ReadFirstLine reads the first line of a text file as a string
func ReadFirstLine(filepath string) (string, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	firstLine := readFirstLineBeyondComments(sc)

	return firstLine, nil
}

func readFirstLineBeyondComments(sc *bufio.Scanner) string {
	inBlockComment := false
	for sc.Scan() {
		line := sc.Text()
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
		return line
	}
	return ""
}

// RemoveExtension returns a path without the extension
func RemoveExtension(relativePath string) string {
	ext := path.Ext(relativePath)
	if ext != "" {
		return relativePath[:len(relativePath)-len(ext)]
	}
	return relativePath
}

// ReadContents reads the entire contents of a text file as a string
func ReadContents(filepath string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return trimByteOrderMark(string(bytes)), nil
}

func trimByteOrderMark(s string) string {
	if len(s) > 3 &&
		s[0] == 0xef &&
		s[1] == 0xbb &&
		s[2] == 0xbf { // byte-order mark
		return s[3:]
	}

	return s
}
