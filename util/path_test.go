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
