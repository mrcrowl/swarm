package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

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

// StringToLines breaks a string in a list of strings, one for each line
func StringToLines(s string) []string {
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		n := len(line) - 1
		if n >= 0 {
			if line[n:] == "\r" {
				lines[i] = line[:n]
			}
		}
	}

	return lines
}

// CountLines returns the number of lines in a string
func CountLines(s string) (int, error) {
	stringReader := strings.NewReader(s)
	reader := bufio.NewReader(stringReader)
	var count int
	for {
		_, isPrefix, err := reader.ReadLine()

		if !isPrefix {
			count++
		}

		if err == io.EOF {
			return count - 1, nil
		}

		if err != nil {
			return count, err
		}
	}

	return count, nil
}
