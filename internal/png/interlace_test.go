package png

import "testing"

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
