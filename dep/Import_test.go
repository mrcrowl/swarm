package dep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathConsistency(t *testing.T) {
	var sut *Import
	sut = newImport("tslib")
	assert.Equal(t, "tslib", sut.path())

	sut = newImport("some/root/relative/path.js")
	assert.Equal(t, "some/root/relative/path.js", sut.path())

	sut = newImport("./one/abc.ts")
	assert.Equal(t, "./one/abc.ts", sut.path())

	sut = newImport("../one/abc.ts")
	assert.Equal(t, "../one/abc.ts", sut.path())
}

func TestDirective(t *testing.T) {
	var sut *Import
	sut = newImport("./search-results.mobile.html#?Config|Config.MOBILE_RELEASE")
	assert.Equal(t, "./search-results.mobile.html", sut.path())
	assert.Equal(t, sut.directive, "#?Config|Config.MOBILE_RELEASE")
}

func TestSolo(t *testing.T) {
	var sut = newImport("tslib")
	assert.True(t, sut.isSolo)
}

func TestRootRelativeImport(t *testing.T) {
	var sut = newImport("app/src/ep/App")
	assert.True(t, sut.isRooted)
	assert.False(t, sut.isSolo)
}

func TestNotRootRelativeImportParent(t *testing.T) {
	var sut = newImport("../../ep/App")
	assert.False(t, sut.isRooted)
}

func TestNotRootRelativeImportSelf(t *testing.T) {
	var sut = newImport("./App")
	assert.False(t, sut.isRooted)
}
