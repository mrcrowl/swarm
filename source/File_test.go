package source

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var temppath string

func setup() {
	temppath, _ = ioutil.TempDir("", "File_test")
}

func teardown() {
	os.RemoveAll(temppath)
}

func getSampleFile(id string, ext string, contents string) *File {
	absoluteFilepath := filepath.Join(temppath, "blah"+ext)
	ioutil.WriteFile(absoluteFilepath, []byte(contents), os.ModePerm)
	return newFile(id, absoluteFilepath)
}

func TestNoImmediateLoad(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".js", "")
	assert.False(t, f.Loaded())
	teardown()
}

func TestLoadNonExistentFile(t *testing.T) {
	f := newFile("blah", "c:\\asldkfhjaksjdfh.js")
	f.EnsureLoaded()
	assert.True(t, f.Loaded())
	assert.IsType(t, &FailedFileContents{}, f.contents)
	assert.Nil(t, f.BundleBody())
}

func TestEnsureLoadedPlainJS(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".js", "alert(\"hi\")")
	f.EnsureLoaded()
	assert.True(t, f.Loaded())
	assert.IsType(t, &JSFileContents{}, f.contents)
	assert.Len(t, f.BundleBody(), 1)
	teardown()
}

func TestEnsureLoadedSystemJS(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".js", `System.register([], function(export, require) { 
		alert(\"hi\");
	}`)
	f.EnsureLoaded()
	assert.True(t, f.Loaded())
	assert.IsType(t, &JSFileContents{}, f.contents)
	assert.Len(t, f.BundleBody(), 3)
	assert.True(t, f.contents.(*JSFileContents).isSystemJS)
	teardown()
}

func TestEnsureLoadedCSS(t *testing.T) {
	setup()
	f := getSampleFile("abcd", ".css", "body { background: green }")
	f.EnsureLoaded()
	assert.True(t, f.Loaded())
	assert.IsType(t, &StringFileContents{}, f.contents)
	assert.Len(t, f.BundleBody(), 13)
	teardown()
}
