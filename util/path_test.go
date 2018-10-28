package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveExtension(t *testing.T) {
	sansExtension := RemoveExtension("/some/path/name.js")
	assert.Equal(t, "/some/path/name", sansExtension)
}

func TestRemoveExtensionWhenNone(t *testing.T) {
	sansExtension := RemoveExtension("/some/path/name")
	assert.Equal(t, "/some/path/name", sansExtension)
}

func TestMimeTypeFromFilename(t *testing.T) {
	cases := map[string]struct {
		filename string
		expected string
	}{
		".js": {
			filename: "blah.js",
			expected: "application/javascript",
		},
		".html": {
			filename: "blah.html",
			expected: "text/html; charset=utf-8",
		},
		".css": {
			filename: "blah.css",
			expected: "text/css; charset=utf-8",
		},
		"???": {
			filename: "akldfoiasudyfiun234",
			expected: "text/plain; charset=utf-8",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := MimeTypeFromFilename(tc.filename)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
