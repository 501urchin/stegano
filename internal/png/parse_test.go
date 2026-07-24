package png

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

var png = []byte{
	0x89, 0x50, 0x4E, 0x47,
	0x0D, 0x0A, 0x1A, 0x0A,

	0x00, 0x00, 0x00, 0x0D,
	0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x01,
	0x08,
	0x06,
	0x00,
	0x00,
	0x00,
	31, 21, 196, 137,

	0x00, 0x00, 0x00, 0x0D,
	0x49, 0x44, 0x41, 0x54,
	0x78, 0x9C,
	0x63, 0x60, 0x60, 0x60,
	0xF8, 0x0F, 0x00, 0x01,
	0x04, 0x01, 0x00,
	95, 229, 195, 75,

	0x00, 0x00, 0x00, 0x00,
	0x49, 0x45, 0x4E, 0x44,
	0xAE, 0x42, 0x60, 0x82,
}

func TestValidateChunk(t *testing.T) {
	source := bytes.NewReader([]byte{1, 2})
	c := PngChunk{
		Length:       2,
		Type:         []byte("IDAT"),
		DataStartIdx: 0,
		CRC:          2347691479,
	}

	err := validateChunk(c, source)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeChunks(t *testing.T) {
	tFilePath := filepath.Join(t.TempDir() + "file.png")
	file, err := os.Create(tFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(png)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Sync()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("valid decode", func(t *testing.T) {
		chunks, err := DecodeChunks(file)
		if err != nil {
			t.Fatal(err)
		}

		if lc := len(chunks); lc != 3 {
			t.Fatalf("failed to return correct number of chunks: want %d, got %d", 3, lc)
		}

		if !slices.Equal(chunks[0].Type, []byte{0x49, 0x48, 0x44, 0x52}) {
			t.Error("chunk index 0 doesnt match expected chunkt type")
		}

		if !slices.Equal(chunks[1].Type, []byte{0x49, 0x44, 0x41, 0x54}) {
			t.Error("chunk index 1 doesnt match expected chunkt type")
		}

		if !slices.Equal(chunks[2].Type, []byte{0x49, 0x45, 0x4E, 0x44}) {
			t.Error("chunk index 2 doesnt match expected chunkt type")
		}

	})

}

func BenchmarkDecodeChunks(b *testing.B) {
	file, err := os.Open("image.png")
	if err != nil {
		panic(err)
	}

	b.Run("custom", func(b *testing.B) {
		for b.Loop() {
			_, _ = file.Seek(0, io.SeekStart)
			_, err = DecodeChunks(file)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

}
