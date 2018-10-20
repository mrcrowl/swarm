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

func TestOutputImports(t *testing.T) {
	var sut = newImportQueue()
	var abcd = newImport("abcd")
	var efgh = newImport("efgh/xyz.js")
	var ijkl = newImport("ijkl/mnop.js")
	sut.push(abcd)
	sut.push(efgh)
	sut.push(ijkl)
	sut.push(abcd)
	imports := sut.OutputImports()
	assert.ElementsMatch(t, []*Import{abcd, ijkl, efgh}, imports)
}
