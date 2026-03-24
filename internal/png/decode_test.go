package png

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func makePng(path string, width, height int) (err error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := range height {
		for x := range width {
			c := color.RGBA{
				R: uint8(x * 25),
				G: uint8(y * 25),
				B: 0,
				A: 255,
			}
			img.Set(x, y, c)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}

func encode(newFileName, oldFileName string, chunks []PngChunk) (err error) {
	newFile, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	oldFile, err := os.Open(oldFileName)
	if err != nil {
		return
	}
	defer oldFile.Close()

	_, err = newFile.Write(pngSignature)
	if err != nil {
		return err
	}

	uint32Buf := make([]byte, 4)
	for i := range chunks {
		binary.BigEndian.PutUint32(uint32Buf, chunks[i].Length)
		_, err = newFile.Write(uint32Buf)
		if err != nil {
			return err
		}

		_, err = newFile.Write(chunks[i].Type)
		if err != nil {
			return err
		}

		_, err = oldFile.Seek(int64(chunks[i].DataStartIdx), io.SeekStart)
		if err != nil {
			return
		}

		if written, cErr := io.CopyN(newFile, oldFile, int64(chunks[i].Length)); cErr != nil || written != int64(chunks[i].Length) {
			return fmt.Errorf("failed to copy correct amount of bytes")
		}

		binary.BigEndian.PutUint32(uint32Buf, chunks[i].CRC)
		_, err = newFile.Write(uint32Buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestDecodePNG(t *testing.T) {
	tFilePath := filepath.Join(t.TempDir(), "test.png")
	err := makePng(tFilePath, 5, 5)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(tFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	t.Run("valid decode", func(t *testing.T) {
		chunks, err := DecodePNG(file)
		if err != nil {
			t.Fatal(err)
		}
		nfileName := filepath.Join(t.TempDir(), "nf.png")

		err = encode(nfileName, tFilePath, chunks)
		if err != nil {
			t.Fatal(err)
		}

		ob, err := os.ReadFile(tFilePath)
		if err != nil {
			t.Fatal(err)
		}
		nb, err := os.ReadFile(nfileName)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(ob, nb) {
			t.Error("failed to decode the correct chunks from png")
		}

	})

}

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
