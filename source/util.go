package source

import (
	"bufio"
	"encoding/json"
	"strings"
)

func stringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func jsonEncodeString(s string) string {
	if b, err := json.Marshal(s); err == nil {
		return string(b)
	}

	return "ERROR encoding"
}
