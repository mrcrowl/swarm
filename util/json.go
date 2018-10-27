package util

import (
	"bytes"
	"encoding/json"
)

// JSONEncodeString converts a string into the equivalent JSON string
func JSONEncodeString(s string) string {
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
