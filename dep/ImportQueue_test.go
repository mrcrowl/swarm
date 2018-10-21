package dep

import (
	"swarm/source"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeenImport(t *testing.T) {
	sut := newImportQueue()
	sut.pushPath("abcd")
	assert.True(t, sut.seen("abcd"))
	assert.False(t, sut.seen("efgh"))
}

func TestUniqueDependency(t *testing.T) {
	sut := newImportQueue()
	sut.pushPath("abcd")
	assert.Equal(t, 1, sut.count())
	sut.pushPath("abcd")
	assert.Equal(t, 1, sut.count())
}

func TestOutputImports(t *testing.T) {
	sut := newImportQueue()
	abcd := source.NewImport("abcd")
	efgh := source.NewImport("efgh/xyz.js")
	ijkl := source.NewImport("ijkl/mnop.js")
	sut.push(abcd)
	sut.push(efgh)
	sut.push(ijkl)
	sut.push(abcd)
	imports := sut.outputImports()
	assert.ElementsMatch(t, []*source.Import{abcd, ijkl, efgh}, imports)
}
