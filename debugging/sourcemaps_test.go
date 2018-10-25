package debugging

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const someJSON = `{
    "version": 3,
    "file": "Third.js",
    "sourceRoot": "",
    "sources": [
        "Third.ts"
    ],
    "names": [],
    "mappings": ";;;;;;;;;YAGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"
}`

func TestParseSourceMap(t *testing.T) {
	value, err := parseSourceMapJSON(someJSON)
	assert.Nil(t, err)
	assert.Equal(t, 3, value.Version, "Version")
	assert.NotEmpty(t, value.File, "File")
	assert.Empty(t, value.SourceRoot, "SourceRoot")
	assert.Len(t, value.Sources, 1, "Sources")
	assert.Len(t, value.Names, 0, "Name")
	assert.NotEmpty(t, value.Mappings, "Mappings")

	parsed := parseMappings(value.Mappings)
	assert.Len(t, parsed, 19)
}

func TestParseMapsString(t *testing.T) {
	mappings := ";;;;AAAA;KAAK;;;;"
	expected := []*line{
		nil,
		nil,
		nil,
		nil,
		&line{
			segments: [][]int{
				[]int{0, 0, 0, 0},
			},
		},
		&line{
			segments: [][]int{
				[]int{5, 0, 0, 5},
			},
		},
		nil,
		nil,
		nil,
		nil,
	}

	actual := parseMappings(mappings)
	assert.Equal(t, len(expected), len(actual), "The # of lines return from parseMaps(...) did not match")
	equal := assert.ObjectsAreEqual(expected, actual)
	if !equal {
		for _, line := range actual {
			log.Printf("%#v\n", line)
		}
		assert.True(t, equal, "The expected result of parseMaps(...) did not match the actual result.")
	}
}
