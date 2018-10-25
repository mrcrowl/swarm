package devtools

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// SourceMap represents the JSON structure of a source map in .map file
type SourceMap struct {
	Version    int      `json:"version"`
	File       string   `json:"file"`
	SourceRoot string   `json:"sourceRoot"`
	Sources    []string `json:"sources"`
	Names      []string `json:"names"`
	Mappings   string   `json:"mappings"`
}

type line struct {
	segments [][]int
}

/*
YAGC = [12,0,3,1]
[
	12, // generated COLUMN (reset with each line, relative within same line)
	0,  // source FILE index (relative to last, except for first) <-- ONLY THING THAT NEEDS TO CHANGE
	4,  // source LINE index (relative to last, except for first)
	1,  // source COLUMN index (relative to last, except for first)
]

 ";;;;AAAA;KAAK;;;;"

*/

func nextNonSeparator(maps string, startPos int) int {
	n := len(maps)
	for i := startPos; i < n; i++ {
		c := maps[i]
		if c != ';' && c != ',' {
			return i
		}
	}
	return -1
}

func nextSeparatorOrEOF(maps string, startPos int) int {
	n := len(maps)
	for i := startPos; i < n; i++ {
		c := maps[i]
		if c == ';' || c == ',' {
			return i
		}
	}
	return n
}

func findFirstVLQ(maps string) (start int, end int) {
	start = nextNonSeparator(maps, 0)
	if start == -1 {
		return -1, -1
	}
	end = nextSeparatorOrEOF(maps, start+1)
	return
}

type vlqReplaceFn func([]int) []int

func replaceFirstVLQ(maps string, replaceFn vlqReplaceFn) string {
	start, end := findFirstVLQ(maps)
	if start < 0 || end < 0 {
		return maps
	}

	before := maps[:start]
	after := maps[end:]
	vlq := maps[start:end]
	values := Decode(vlq)
	replacementValues := replaceFn(values)
	replacementVlq := Encode(replacementValues)
	return before + replacementVlq + after
}

func parseMappings(maps string) []*line {
	lineStrings := strings.Split(maps, ";")
	lines := make([]*line, len(lineStrings))
	for i, lineString := range lineStrings {
		lines[i] = parseLineString(lineString)
	}
	return lines
}

func parseLineString(lineString string) *line {
	if lineString == "" {
		return nil
	}
	segmentStrings := strings.Split(lineString, ",")
	segments := make([][]int, len(segmentStrings))
	for i, segmentString := range segmentStrings {
		segments[i] = Decode(segmentString)
	}
	return &line{segments}
}

func parseSourceMapJSON(sourceMapJSON string) (*SourceMap, error) {
	var sm *SourceMap
	err := json.Unmarshal([]byte(sourceMapJSON), &sm)
	if err != nil {
		return nil, errors.New("Invalid JSON in source map: " + err.Error())
	}
	return sm, nil
}

const base64Map = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="

func byteToInt(b byte) int {
	switch {
	case b >= 'A' && b <= 'Z':
		return int(b - 'A')
	case b >= 'a' && b <= 'z':
		return int(b - 'a' + 26)
	case b >= '0' && b <= '9':
		return int(b - '0' + 52)
	case b == '+':
		return 62
	case b == '/':
		return 63
	case b == '=':
		return 64
	default:
		panic(fmt.Sprintf("byteToInt received byte out of range: %c", b))
	}
}

func intToByte(i int) byte {
	if i >= 0 && i <= 64 {
		return base64Map[i]
	}

	panic(fmt.Sprintf("intToByte received int out of range: %d", i))
}

// Decode decodes a base-64 VLQ string to a list of integers
func Decode(s string) []int {
	result := make([]int, 0, 8)
	shift := uint(0)
	value := 0

	for _, b := range s {
		integer := byteToInt(byte(b))

		hasContinuationBit := (integer & 32) > 0

		integer &= 31
		value += integer << shift

		if hasContinuationBit {
			shift += 5
		} else {
			shouldNegate := (value & 1) > 0
			value >>= 1

			if shouldNegate {
				result = append(result, -value)
			} else {
				result = append(result, value)
			}

			// reset
			value = 0
			shift = 0
		}
	}

	return result
}

// Encode encodes a list of numbers to a VLQ string
func Encode(values []int) string {
	result := make([]byte, 0, 16)
	for _, n := range values {
		result = append(result, encodeInteger(n)...)
	}
	return string(result)
}

func encodeInteger(n int) []byte {
	result := make([]byte, 0, 8)

	if n < 0 {
		n = (-n << 1) | 1
	} else {
		n <<= 1
	}

	for {
		clamped := n & 31
		n >>= 5

		if n > 0 {
			clamped |= 32
		}

		result = append(result, intToByte(clamped))

		if n <= 0 {
			break
		}
	}

	return result
}
