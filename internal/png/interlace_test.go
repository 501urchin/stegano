package png

import (
	"slices"
	"testing"
)

func TestBytePerRow(t *testing.T) {
	tests := []struct {
		name      string
		width     int
		interlace byte
		colorType byte
		bitDepth  byte
		expected  int
	}{
		{
			name:      "grayscale 8 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 0,
			bitDepth:  8,
			expected:  11,
		},
		{
			name:      "rgb 8 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 2,
			bitDepth:  8,
			expected:  31,
		},
		{
			name:      "rgba 8 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 6,
			bitDepth:  8,
			expected:  41,
		},
		{
			name:      "grayscale alpha 8 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 4,
			bitDepth:  8,
			expected:  21,
		},
		{
			name:      "indexed color 8 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 3,
			bitDepth:  8,
			expected:  11,
		},
		{
			name:      "indexed color 4 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 3,
			bitDepth:  4,
			expected:  6,
		},
		{
			name:      "indexed color 1 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 3,
			bitDepth:  1,
			expected:  3,
		},
		{
			name:      "rgba 16 bit no interlace",
			width:     10,
			interlace: 0,
			colorType: 6,
			bitDepth:  16,
			expected:  81,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bpr, err := BytePerRow(uint32(test.width), test.interlace, test.colorType, test.bitDepth)
			if err != nil {
				t.Error(err)
			}

			if bpr != test.expected {
				t.Errorf("result (%d) doesnt match expected value (%d)", bpr, test.expected)
			}
		})
	}
}

func TestReverseNoneFilter(t *testing.T) {
	tests := []struct {
		name     string
		current  []byte
		previous []byte
		bpp      int
		expected []byte
	}{

		{
			name:     "none filter",
			current:  []byte{0, 10, 20, 30, 40},
			previous: []byte{0, 5, 10, 15, 20},
			bpp:      1,
			expected: []byte{0, 10, 20, 30, 40},
		},

		{
			name:     "sub filter first row",
			current:  []byte{1, 10, 20, 30, 40},
			previous: nil,
			bpp:      1,
			expected: []byte{1, 10, 30, 60, 100},
		},

		{
			name:     "sub filter rgb pixels",
			current:  []byte{1, 10, 20, 30, 5, 10, 15},
			previous: nil,
			bpp:      3,
			expected: []byte{1, 10, 20, 30, 15, 30, 45},
		},

		{
			name:     "up filter first row",
			current:  []byte{2, 10, 20, 30},
			previous: nil,
			bpp:      1,
			expected: []byte{2, 10, 20, 30},
		},

		{
			name:     "up filter",
			current:  []byte{2, 10, 20, 30},
			previous: []byte{2, 5, 10, 20},
			bpp:      1,
			expected: []byte{2, 15, 30, 50},
		},

		{
			name:     "average filter first row",
			current:  []byte{3, 10, 20, 30},
			previous: nil,
			bpp:      1,
			expected: []byte{3, 10, 25, 42},
		},

		{
			name:     "average filter",
			current:  []byte{3, 10, 20, 30},
			previous: []byte{3, 20, 30, 40},
			bpp:      1,
			expected: []byte{3, 20, 45, 72},
		},

		{
			name:     "paeth filter first row",
			current:  []byte{4, 10, 20, 30},
			previous: nil,
			bpp:      1,
			expected: []byte{4, 10, 30, 60},
		},

		{
			name:     "paeth filter",
			current:  []byte{4, 10, 10, 10},
			previous: []byte{4, 20, 20, 20},
			bpp:      1,
			expected: []byte{4, 30, 40, 50},
		},

		{
			name:     "byte overflow",
			current:  []byte{2, 250},
			previous: []byte{2, 20},
			bpp:      1,
			expected: []byte{2, 14},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ReverseNoneFilter(test.current, test.previous, test.bpp)

			if !slices.Equal(result, test.expected) {
				t.Errorf("result doesnt match expected. got %v, want %v", result, test.expected)
			}
		})
	}
}
