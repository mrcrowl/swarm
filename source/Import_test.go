package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathConsistency(t *testing.T) {
	var sut *Import
	sut = NewImport("tslib")
	assert.Equal(t, "tslib", sut.Path())

	sut = NewImport("some/root/relative/path.js")
	assert.Equal(t, "some/root/relative/path.js", sut.Path())

	sut = NewImport("./one/abc.ts")
	assert.Equal(t, "./one/abc.ts", sut.Path())

	sut = NewImport("../one/abc.ts")
	assert.Equal(t, "../one/abc.ts", sut.Path())
}

func TestDirective(t *testing.T) {
	var sut *Import
	sut = NewImport("./search-results.mobile.html#?Config|Config.MOBILE_RELEASE")
	assert.Equal(t, "./search-results.mobile.html", sut.Path())
	assert.Equal(t, sut.Directive, "#?Config|Config.MOBILE_RELEASE")
}

func TestSolo(t *testing.T) {
	var sut = NewImport("tslib")
	assert.True(t, sut.IsSolo)
}

func TestRootRelativeImport(t *testing.T) {
	var sut = NewImport("app/src/ep/App")
	assert.True(t, sut.IsRooted)
	assert.False(t, sut.IsSolo)
}

func TestNotRootRelativeImportParent(t *testing.T) {
	var sut = NewImport("../../ep/App")
	assert.False(t, sut.IsRooted)
}

func TestNotRootRelativeImportSelf(t *testing.T) {
	var sut = NewImport("./App")
	assert.False(t, sut.IsRooted)
}
