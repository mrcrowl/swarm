package bundle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			actual := Decode(tc.vlq)
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
			actual := Encode(tc.nums)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
