package devtools

import (
	"fmt"
	"strings"
	"swarm/source"
)

// PlayMappings loops through the mappings to calculate a "delta" that occurs
// by applying "the rules".
func (smap *sourceMap) PlayMappings() *source.MapPlayback {
	playback := smap.mapping.Playback()
	if playback == nil {
		var segDelta source.Segment
		lines := parseMappings(smap.mapping.Mappings())
		for _, line := range lines {
			if line != nil {
				segDelta.GeneratedColumn = 0
				for _, seg := range line.segments {
					segDelta = segDelta.Add(*seg)
					// fmt.Printf("[%d,%d](#%d)=>[%d,%d] |", segDelta.sourceLine, segDelta.sourceColumn,
					// segDelta.sourceFile, generatedLine, segDelta.generatedColumn)
				}
			}
			// fmt.Println()
		}
		playback = &source.MapPlayback{LineCount: len(lines), SegmentDelta: segDelta}
		smap.mapping.CachePlayback(playback)
	}
	return playback
}

// OffsetMappings replaces the source file index of the first
// VLQ in the Mappings field of this smap.  This is used for concatenating multiple source maps together.
// See: https://sourcemaps.info/spec.html
//      http://www.murzwin.com/base64vlq.html (WARNING: the ability to "play" source maps, near the bottom of this page is incorrect for this site!)
func (smap *sourceMap) OffsetMappings(segDelta source.Segment) string {
	offsetMappings := replaceFirstVLQ(smap.mapping.Mappings(), func(seg source.Segment) source.Segment {
		adjustedSeg := segDelta.AdjustForSource()
		resetSeg := seg.Add(adjustedSeg)
		return resetSeg
	})
	return offsetMappings
}

type sourceMap struct {
	spacerLines   int
	fileLineCount int
	mapping       *source.Mapping
}

// func (smap *sourceMap) contents() string {
// 	return smap.mapping.Contents()
// }

func (smap *sourceMap) path() string {
	return smap.mapping.RelativePath()
}

type line struct {
	segments []*source.Segment
}

/*
	Mappings := ";;;;YAGC;KAAK;;;;"

	YAGC = [12,0,3,1]
	[
		12, // generated COLUMN (reset with each line, relative within same line)
		0,  // source FILE index (relative to last, except for first) <-- ONLY THING THAT NEEDS TO CHANGE
		4,  // source LINE index (relative to last, except for first)
		1,  // source COLUMN index (relative to last, except for first)
	]
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

type vlqReplaceFn func(source.Segment) source.Segment

func replaceFirstVLQ(mappings string, replaceFn vlqReplaceFn) string {
	start, end := findFirstVLQ(mappings)
	if start < 0 || end < 0 {
		return mappings
	}

	before := mappings[:start]
	after := mappings[end:]
	vlq := mappings[start:end]
	values := decodeSegment(vlq)
	replacementValues := replaceFn(values)
	replacementVlq := encodeSegment(replacementValues)
	return before + replacementVlq + after
}

func parseMappings(mappings string) []*line {
	lineStrings := strings.Split(mappings, ";")
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
	segments := make([]*source.Segment, len(segmentStrings))
	for i, segmentString := range segmentStrings {
		seg := decodeSegment(segmentString)
		segments[i] = &seg
	}
	return &line{segments}
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

// decode decodes a base-64 VLQ string to a strongly-typed segment
func decodeSegment(s string) source.Segment {
	values := decode(s)
	if len(values) >= 4 {
		return source.Segment{
			GeneratedColumn: values[0],
			SourceFile:      values[1],
			SourceLine:      values[2],
			SourceColumn:    values[3],
		}
	}
	panic(fmt.Sprintf("Encountered decode result with fewer than 4 values: %#v", values))
}

// decode decodes a base-64 VLQ string to a list of integers
func decode(s string) []int {
	result := make([]int, 0, 4)
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

// encode encodes a list of numbers to a VLQ string

func encodeSegment(seg source.Segment) string {
	values := []int{seg.GeneratedColumn, seg.SourceFile, seg.SourceLine, seg.SourceColumn}
	return encode(values)
}

// encode encodes a list of numbers to a VLQ string
func encode(values []int) string {
	result := make([]byte, 0, 8)
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
