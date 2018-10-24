package bundle

import (
	"fmt"
)

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

/*
function encodeInteger ( num: number ): string {
	let result = '';

	if ( num < 0 ) {
		num = ( -num << 1 ) | 1;
	} else {
		num <<= 1;
	}

	do {
		let clamped = num & 31;
		num >>= 5;

		if ( num > 0 ) {
			clamped |= 32;
		}

		result += integerToChar[ clamped ];
	} while ( num > 0 );

	return result;
}
*/
