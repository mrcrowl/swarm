package devtools

import (
	"log"
	"github.com/mrcrowl/swarm/source"
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

const firstMappingsNoChange = ";;;;;;;;;YAGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"
const firstMappingsPlusOne1 = ";;;;;;;;;YCGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"
const firstMappingsPlus1506 = ";;;;;;;;;Yk+CGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"
const firstMappingsMinus369 = ";;;;;;;;;YjXGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"

const thirdJSON = `{
    "version": 3,
    "file": "Third.js",
    "sourceRoot": "",
    "sources": [
        "Third.ts"
    ],
    "names": [],
    "mappings": ";;;;;;;;;YAGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"
}`

const mapping1 = `;;;;YAGA;gBAAA;gBAIA,CAAC`
const mapping2 = `;;;;;;YACA;gBAAA;gBAKA,CAAC;`
const combinedMappings = `;;;;YAGA;gBAAA;gBAIA,CAAC;;;;;;YCCA;gBAAA;gBAKA,CAAC;`

func TestPlayMappings(t *testing.T) {
	cases := map[string]struct {
		mappings          string
		expected          source.Segment
		expectedLineCount int
	}{
		"zero":          {expectedLineCount: 1, mappings: "AAAA", expected: source.Segment{GeneratedColumn: 0, SourceFile: 0, SourceLine: 0, SourceColumn: 0}},
		"1234":          {expectedLineCount: 1, mappings: "ACEG", expected: source.Segment{GeneratedColumn: 0, SourceFile: 1, SourceLine: 2, SourceColumn: 3}},
		"1234;;":        {expectedLineCount: 1, mappings: "ACEG", expected: source.Segment{GeneratedColumn: 0, SourceFile: 1, SourceLine: 2, SourceColumn: 3}},
		"AAAA":          {expectedLineCount: 1, mappings: "AAAA", expected: source.Segment{GeneratedColumn: 0, SourceFile: 0, SourceLine: 0, SourceColumn: 0}},
		"MED":           {expectedLineCount: 4, mappings: "AAAA;BBBB;CCCC,ACCC,ABBB,XYZA;ADDD", expected: source.Segment{GeneratedColumn: 0, SourceFile: 13, SourceLine: -11, SourceColumn: 1}},
		"LONG":          {expectedLineCount: 19, expected: source.Segment{GeneratedColumn: 9, SourceFile: 0, SourceLine: 7, SourceColumn: 3}, mappings: ";;;;;;;;;YAGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,qCAAqC,CAAA;gBAChD,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAC,CAAC"},
		"Config.js.map": {expectedLineCount: 121, expected: source.Segment{GeneratedColumn: 9, SourceFile: 0, SourceLine: 173, SourceColumn: 1}, mappings: "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"},
		// "First.js.map":  {expectedLineCount: 3, expected: source.Segment{GeneratedColumn:9,SourceFile: 0, SourceLine:7, SourceColumn:2}, mappings: ";;;;;;;;;;YAGA;gBAAA;gBAIA,CAAC;gBAHiB,QAAE,GAAhB;oBACI,OAAO,eAAa,eAAM,CAAC,EAAE,EAAI,CAAA;gBACrC,CAAC;gBACL,YAAC;YAAD,CAAC,AAJD,IAIC;;QAAA,CAAC"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			config := &source.MapConfig{Mappings: tc.mappings}
			mapping := source.NewMappingForTesting(config)
			smap := &sourceMap{mapping: mapping}
			playback := smap.PlayMappings()
			assert.Equal(t, tc.expectedLineCount, playback.LineCount)
			assert.Equal(t, tc.expected, playback.SegmentDelta)
		})
	}
}

// func TestOffsetMappingsSourceFileIndex(t *testing.T) {
// 	cases := map[string]struct {
// 		json      string
// 		fileIndex int
// 		expected  string
// 	}{
// 		"no-change": {
// 			json:      firstJSON,
// 			fileIndex: 0,
// 			expected:  firstMappingsNoChange,
// 		},
// 		"increase-by-1": {
// 			json:      firstJSON,
// 			fileIndex: 1,
// 			expected:  firstMappingsPlusOne1,
// 		},
// 		"increase-by-1506": {
// 			json:      firstJSON,
// 			fileIndex: 1506,
// 			expected:  firstMappingsPlus1506,
// 		},
// 		"decrease-by-369": { // <-- this one is stupid, but meh \_/
// 			json:      firstJSON,
// 			fileIndex: -369,
// 			expected:  firstMappingsMinus369,
// 		},
// 	}
// 	for name, tc := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			smap, err := ParseSourceMapJSON(tc.json)
// 			assert.Nil(t, err)
// 			actual := smap.OffsetMappingsSourceFileIndex(tc.fileIndex)
// 			assert.Equal(t, tc.expected, actual)
// 		})
// 	}
// }

func TestFindFirstLVQ(t *testing.T) {
	cases := map[string]struct {
		mappings      string
		startPos      int
		expectedStart int
		expectedEnd   int
	}{
		"one":   {mappings: ";;;;;AAAA", startPos: 0, expectedStart: 5, expectedEnd: 9},
		"two":   {mappings: ";;AAAA;;;AZQA;bGAFA;", startPos: 2, expectedStart: 2, expectedEnd: 6},
		"none":  {mappings: ";;;;;;;;;;;;", startPos: 0, expectedStart: -1, expectedEnd: -1},
		"start": {mappings: "AAAA;;;;;", startPos: 9, expectedStart: 0, expectedEnd: 4},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			start, end := findFirstVLQ(tc.mappings)
			assert.Equal(t, tc.expectedStart, start)
			assert.Equal(t, tc.expectedEnd, end)
		})
	}
}

func TestNextNonSeparator(t *testing.T) {
	cases := map[string]struct {
		mappings string
		startPos int
		expected int
	}{
		"one":   {mappings: ";;;;;AAAA", startPos: 0, expected: 5},
		"two":   {mappings: ";;;;;AAAA", startPos: 2, expected: 5},
		"start": {mappings: "AAAC;AAAD;ZZZA", startPos: 0, expected: 0},
		"eof-1": {mappings: "AAAC;AAAD;;;;;", startPos: 9, expected: -1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := nextNonSeparator(tc.mappings, tc.startPos)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestNextSepartorOrEOF(t *testing.T) {
	cases := map[string]struct {
		mappings string
		startPos int
		expected int
	}{
		"start":  {mappings: ";;;;;AAAA", startPos: 0, expected: 0},
		"eof":    {mappings: ";;;;;AAAA", startPos: 5, expected: 9},
		"second": {mappings: "A;B", startPos: 0, expected: 1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := nextSeparatorOrEOF(tc.mappings, tc.startPos)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestReplaceFirstVLQ(t *testing.T) {
	cases := map[string]struct {
		mappings      string
		replacementFn vlqReplaceFn
		expected      string
	}{
		"one": {
			mappings: "YCCA",
			replacementFn: func(seg source.Segment) source.Segment {
				seg.GeneratedColumn++
				return seg
			},
			expected: "aCCA",
		},
		"upndown": {
			mappings: "AAAA",
			replacementFn: func(seg source.Segment) source.Segment {
				seg.GeneratedColumn++
				seg.SourceFile--
				seg.SourceLine++
				seg.SourceColumn--
				return seg
			},
			expected: "CDCD",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := replaceFirstVLQ(tc.mappings, tc.replacementFn)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

/*
YAGC = [12,0,3,1]
[
	12, // generated COLUMN (reset with each line, relative within same line)
	0,  // source FILE index (relative to last, except for first) <-- ONLY THING THAT NEEDS TO CHANGE
	4,  // source LINE index (relative to last, except for first)
	1,  // source COLUMN index (relative to last, except for first)
]
*/

func TestParseMapsString(t *testing.T) {
	mappings := ";;;;AAAA;KAAK;;;;"
	expected := []*line{
		nil,
		nil,
		nil,
		nil,
		&line{
			segments: []*source.Segment{
				&source.Segment{GeneratedColumn: 0, SourceFile: 0, SourceLine: 0, SourceColumn: 0},
			},
		},
		&line{
			segments: []*source.Segment{
				&source.Segment{GeneratedColumn: 5, SourceFile: 0, SourceLine: 0, SourceColumn: 5},
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

func TestDecode(t *testing.T) {
	cases := map[string]struct {
		vlq      string
		expected []int
	}{
		"AAAC": {
			vlq:      "AAAC",
			expected: []int{0, 0, 0, 1},
		},
		"ADAA": {
			vlq:      "ADAA",
			expected: []int{0, -1, 0, 0},
		},
		"AAgBC": {
			vlq:      "AAgBC",
			expected: []int{0, 0, 16, 1},
		},
		"KAAK": {
			vlq:      "KAAK",
			expected: []int{5, 0, 0, 5},
		},
		"G9s6a8zns//+": {
			vlq:      "G9s6aAs8BzC",
			expected: []int{3, -439502, 0, 966, -41},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := decode(tc.vlq)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestEncode(t *testing.T) {
	cases := map[string]struct {
		expected string
		nums     []int
	}{
		"AAAC": {
			nums:     []int{0, 0, 0, 1},
			expected: "AAAC",
		},
		"ADAA": {
			nums:     []int{0, -1, 0, 0},
			expected: "ADAA",
		},
		"AAgBC": {
			nums:     []int{0, 0, 16, 1},
			expected: "AAgBC",
		},
		"KAAK": {
			nums:     []int{5, 0, 0, 5},
			expected: "KAAK",
		},
		"G9s6a8zns//+": {
			nums:     []int{3, -439502, 0, 966, -41},
			expected: "G9s6aAs8BzC",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := encode(tc.nums)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
