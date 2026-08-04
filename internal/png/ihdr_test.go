package png

import (
	"bytes"
	"reflect"
	"testing"

	pngerrors "github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

func TestParseIHDRChunk(t *testing.T) {
	ihdrChunk := []byte{
		// Length (13)
		0x00, 0x00, 0x00, 0x0D,

		// Chunk type: IHDR
		0x49, 0x48, 0x44, 0x52,

		// Data
		0x00, 0x00, 0x00, 0x01, // Width = 1
		0x00, 0x00, 0x00, 0x01, // Height = 1
		0x08, // Bit depth = 8
		0x06, // Color type = 6 (RGBA)
		0x00, // Compression = 0
		0x00, // Filter = 0
		0x00, // Interlace = 0

		// CRC (for "IHDR" + data)
		0x1F, 0x15, 0xC4, 0x89,
	}
	expected := IHDRData{
		Width:       1,
		Height:      1,
		BitDepth:    8,
		ColorType:   6,
		Compression: 0,
		Filter:      0,
		Interlace:   0,
	}

	src := bytes.NewReader(ihdrChunk)
	chunk := PngChunk{
		Length:       13,
		Type:         types.IHDR,
		DataStartIdx: 8,
	}

	d, err := ParseIHDRChunk(src, chunk)
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
			name: "invalid color type",
			data: []byte{
				0x00, 0x00, 0x00, 0x0D,
				0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x01,
				0x08,
				0x07,
				0x00,
				0x00,
				0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: pngerrors.ErrInvalidPNGColorType,
		},
		{
			name: "invalid compression method",
			data: []byte{
				0x00, 0x00, 0x00, 0x0D,
				0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x01,
				0x08,
				0x06,
				0x01,
				0x00,
				0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: pngerrors.ErrInvalidPNGCompression,
		},
		{
			name: "invalid filter method",
			data: []byte{
				0x00, 0x00, 0x00, 0x0D,
				0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x01,
				0x08,
				0x06,
				0x00,
				0x01,
				0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: pngerrors.ErrInvalidPNGFilter,
		},
		{
			name: "invalid interlace method",
			data: []byte{
				0x00, 0x00, 0x00, 0x0D,
				0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x01,
				0x08,
				0x06,
				0x00,
				0x00,
				0x02,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: pngerrors.ErrInvalidPNGInterlace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := bytes.NewReader(tt.data)
			chunk := PngChunk{
				Length:       13,
				Type:         types.IHDR,
				DataStartIdx: 8,
			}

			_, err := ParseIHDRChunk(src, chunk)
			if err == nil {
				t.Fatal("expected an error")
			}
			if err != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, err)
			}
		})
	}
}
