package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

func stringToLines(s string) []string {
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

func jsonEncodeString(s string) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(s); err == nil {
		s := buf.String()
		n := len(s) - 1
		if s[n] == '\n' {
			return s[:n]
		}
		return s
	}

	return "ERROR encoding"
}

func countLines(s string) (int, error) {
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
		} else if err != nil {
			return count, err
		}
	}

	return count, nil
}
