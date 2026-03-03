package compression

import (
	"bytes"
	"testing"
)

func TestCompressZSTD(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "compress small data",
			data:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "compress large data",
			data:    make([]byte, 10000), // Large data
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompressZSTD(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressZSTD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("CompressZSTD() returned empty compressed data")
			}
		})
	}
}

func TestDecompressZSTD(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "decompress valid data",
			data:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "decompress empty data",
			data:    []byte(""),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compress data first for decompression testing
			compressedData, err := CompressZSTD(tt.data)
			if err != nil {
				t.Fatalf("CompressZSTD() failed: %v", err)
			}

			// Now test decompression
			got, err := DecompressZSTD(compressedData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecompressZSTD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(got, tt.data) {
				t.Errorf("DecompressZSTD() = %v, want %v", got, tt.data)
			}
		})
	}
}

func TestRoundTripCompressionDecompression(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "round-trip compression and decompression",
			data:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "round-trip empty data",
			data:    []byte(""),
			wantErr: false,
		},
		{
			name:    "round-trip large data",
			data:    make([]byte, 10000), // Large data
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compress
			compressedData, err := CompressZSTD(tt.data)
			if err != nil {
				t.Fatalf("CompressZSTD() failed: %v", err)
			}

			// Decompress
			got, err := DecompressZSTD(compressedData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecompressZSTD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(got, tt.data) {
				t.Errorf("Round-trip data mismatch: got %v, want %v", got, tt.data)
			}
		})
	}
}

func TestDecompressInvalidData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "decompress invalid compressed data",
			data:    []byte("invalid data"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecompressZSTD(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecompressZSTD() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got) > 0 {
				t.Errorf("DecompressZSTD() = %v, want error", got)
			}
		})
	}
}
