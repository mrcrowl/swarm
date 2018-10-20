package dep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeenImport(t *testing.T) {
	var sut = newImportQueue()
	sut.pushPath("abcd")
	assert.True(t, sut.seen("abcd"))
	assert.False(t, sut.seen("efgh"))
}

func TestUniqueDependency(t *testing.T) {
	var sut = newImportQueue()
	sut.pushPath("abcd")
	assert.Equal(t, 1, sut.count())
	sut.pushPath("abcd")
	assert.Equal(t, 1, sut.count())
}
