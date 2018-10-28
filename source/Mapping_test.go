package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const firstJSON = `{
    "version": 3,
    "file": "First.js",
    "sourceRoot": "",
    "sources": [
        "First.ts"
    ],
    "names": [],
    "mappings": ";;;;;;;;;YAGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"
}`

func TestParseSourceMap(t *testing.T) {
	value, err := ParseSourceMapConfig(firstJSON)
	assert.Nil(t, err)
	assert.Equal(t, 3, value.Version, "Version")
	assert.NotEmpty(t, value.File, "File")
	assert.Empty(t, value.SourceRoot, "SourceRoot")
	assert.Len(t, value.Sources, 1, "Sources")
	assert.Len(t, value.Names, 0, "Name")
	assert.NotEmpty(t, value.Mappings, "Mappings")

	// parsed := parseMappings(value.Mappings)
	// assert.Len(t, parsed, 19)
}
