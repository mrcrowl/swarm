package io

import (
	"os"
	"swarm/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var temppath string

func setup() {
	temppath = testutil.CreateTempDirWithPrefix("fileutil")
}

func teardown() {
	os.RemoveAll(temppath)
}

func TestRemoveExtension(t *testing.T) {
	sansExtension := RemoveExtension("/some/path/name.js")
	assert.Equal(t, "/some/path/name", sansExtension)
}

func TestRemoveExtensionWhenNone(t *testing.T) {
	sansExtension := RemoveExtension("/some/path/name")
	assert.Equal(t, "/some/path/name", sansExtension)
}

const randomContents = `random\n\r\n\t\0whatever`

func TestReadContents(t *testing.T) {
	setup()
	filepath := testutil.WriteTextFile(temppath, "TestReadContents", randomContents)
	readContents, err := ReadContents(filepath)
	assert.Nil(t, err)
	assert.Equal(t, randomContents, readContents)
	teardown()
}

func TestReadContentsMissing(t *testing.T) {
	readContents, err := ReadContents("aksjldfhaskjeh98dfjkahf.zjkdhfa")
	assert.NotNil(t, err)
	assert.Equal(t, "", readContents)
}

func TestReadFirstLineMissing(t *testing.T) {
	line, err := ReadFirstLine("aksjldfhaskjeh98dfjkahf.zjkdhfa")
	assert.NotNil(t, err)
	assert.Equal(t, "", line)
}

func TestReadFirstLineEmpty(t *testing.T) {
	setup()
	source := ""
	filepath := testutil.WriteTextFile(temppath, "TestReadFirstLineEmpty", source)
	line, err := ReadFirstLine(filepath)
	assert.Nil(t, err)
	assert.Equal(t, "", line)
	teardown()
}

func TestReadFirstLinePrefix(t *testing.T) {
	setup()
	source := `// comment 1
// comment 2
abcd`
	filepath := testutil.WriteTextFile(temppath, "TestReadFirstLinePrefix", source)
	firstLine, err := ReadFirstLine(filepath)
	assert.Nil(t, err)
	assert.Equal(t, "abcd", firstLine)
	teardown()
}

func TestReadFirstLineBlock(t *testing.T) {
	setup()
	source := `/* comment 1
 comment 2
 comment 3 */ 
abcd`
	filepath := testutil.WriteTextFile(temppath, "TestReadFirstLineBlock", source)
	firstLine, err := ReadFirstLine(filepath)
	assert.Nil(t, err)
	assert.Equal(t, "abcd", firstLine)
	teardown()
}

func TestReadFirstLineMixture(t *testing.T) {
	setup()
	source := `/* comment 1
 comment 2
 comment 3 */ 
// another rand comment 4
abcd`
	filepath := testutil.WriteTextFile(temppath, "TestReadFirstLineMixture", source)
	firstLine, err := ReadFirstLine(filepath)
	assert.Nil(t, err)
	assert.Equal(t, "abcd", firstLine)
	teardown()
}
