package png

import (
	"reflect"
	"testing"
)

func TestParseIHDRChunk(t *testing.T) {
	ihdrRGBA8 := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00}

	expected := PngInfo{
		Width:       1,
		Height:      1,
		BitDepth:    8,
		ColorType:   6,
		Compression: 0,
		Filter:      0,
		Interlace:   0,
	}

	d, err := ParseIHDRChunk(ihdrRGBA8)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(d, expected) {
		t.Errorf("expected %+v, got %+v", expected, d)
	}
}

func TestParseIHDRChunkErrors(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected error
	}{
		{
			name:     "invalid color type",
			data:     []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x07, 0x00, 0x00, 0x00},
			expected: ErrInvalidPNGColorType,
		},
		{
			name:     "invalid compression method",
			data:     []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x01, 0x00, 0x00},
			expected: ErrInvalidPNGCompression,
		},
		{
			name:     "invalid filter method",
			data:     []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x01, 0x00},
			expected: ErrInvalidPNGFilter,
		},
		{
			name:     "invalid interlace method",
			data:     []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x02},
			expected: ErrInvalidPNGInterlace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseIHDRChunk(tt.data)
			if err == nil {
				t.Fatal("expected an error")
			}
			if err != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, err)
			}
		})
	}
}
