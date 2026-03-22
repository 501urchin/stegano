package png

import (
	"image/png"
	"io"
	"os"
	"testing"
)

// func TestDecodePNG(t *testing.T) {
// 	testPngPath := filepath.Join(t.TempDir(), "test.png")

// 	err := os.WriteFile(testPngPath, testPng, 0644)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	file, err := os.Open(testPngPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer file.Close()

// 	t.Run("valid decode", func(t *testing.T) {
// 		chunks, err := DecodePNG(file)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		assert.Equal(t, expectedChunk, chunks)
// 	})

// }

func BenchmarkDecodePNG(b *testing.B) {
	file, err := os.Open("image.png")
	if err != nil {
		panic(err)
	}

	b.Run("custom", func(b *testing.B) {
		for b.Loop() {
			_, _ = file.Seek(0, io.SeekStart)
			_, err = DecodePNG(file)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("image.Image version", func(b *testing.B) {
		for b.Loop() {
			_, _ = file.Seek(0, io.SeekStart)
			_, err = png.Decode(file)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
