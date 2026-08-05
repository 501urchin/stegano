package stegano

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	ierr "github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

// TestCase for Flush
// test if dst image fully matches src
// test flush at middle of image

func TestNewPngEncoder(t *testing.T) {
	srcFileName := t.TempDir() + "src.png"
	srcFile, err := os.Create(srcFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()
	defer os.Remove(srcFileName)

	dstFileName := t.TempDir() + "dst.png"
	dstFile, err := os.Create(dstFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()
	defer os.Remove(dstFileName)

	resetFile := func(f *os.File) (err error) {
		err = f.Truncate(0)
		if err != nil {
			return
		}

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return
		}

		return f.Sync()
	}

	t.Run("invalid bitdepth", func(t *testing.T) {
		for i := 8; i < 256; i++ {
			_, err := NewPngEncoder(srcFile, dstFile, types.BitIndex(i))
			if !errors.Is(err, ierr.ErrInvalidBitDepth) {
				t.Fatalf("failed to return the desired error: want %q but got %q", ierr.ErrInvalidBitDepth, err)
			}
		}
	})

	t.Run("Nil src", func(t *testing.T) {
		_, err = NewPngEncoder(nil, dstFile, 0)
		if !errors.Is(err, ierr.ErrSrcIsNil) {
			t.Errorf("failed to return the desired error: want %q but got %q", ierr.ErrSrcIsNil, err)
		}
	})

	t.Run("Nil dst", func(t *testing.T) {
		_, err = NewPngEncoder(srcFile, nil, 0)
		if !errors.Is(err, ierr.ErrDstIsNil) {
			t.Errorf("failed to return the desired error: want %q but got %q", ierr.ErrDstIsNil, err)
		}
	})

	t.Run("Invalid png (missing header)", func(t *testing.T) {
		_, err = NewPngEncoder(srcFile, dstFile, 0)
		if !errors.Is(err, ierr.ErrFailedToReadSignature) {
			t.Errorf("failed to return the desired error: want %q but got %q", ierr.ErrFailedToReadSignature, err)
		}
	})

	t.Run("Invalid png (invalid header)", func(t *testing.T) {
		_, err = srcFile.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x1A})
		if err != nil {
			t.Fatal(err)
		}
		err = srcFile.Sync()
		if err != nil {
			t.Fatal(err)
		}

		_, err = NewPngEncoder(srcFile, dstFile, 0)
		if !errors.Is(err, ierr.ErrNotPNG) {
			t.Errorf("failed to return the desired error: want %q but got %q", ierr.ErrNotPNG, err)
		}

		err = resetFile(srcFile)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid png (missing ihdr)", func(t *testing.T) {
		_, err = srcFile.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
		if err != nil {
			t.Fatal(err)
		}
		srcFile.Sync()

		_, err = NewPngEncoder(srcFile, dstFile, 0)
		if !errors.Is(err, ierr.ErrMissingIHDR) {
			t.Errorf("failed to return the desired error: want %q but got %q", ierr.ErrMissingIHDR, err)
		}

		resetFile(srcFile)
		resetFile(dstFile)
	})

	t.Run("wrote chunks before IDAT to dst", func(t *testing.T) {
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

		_, err = srcFile.Write(png)
		if err != nil {
			t.Fatal(err)
		}
		srcFile.Sync()

		_, err = NewPngEncoder(srcFile, dstFile, 0)
		if err != nil {
			t.Fatal(err)
		}
		dstFile.Sync()

		_, err = dstFile.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 33)

		_, err = dstFile.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(buf, png[:33]) {
			t.Errorf("failed to write chunks before IDAT to dst. (\nExpected: %q, \nBut got: %q)", png[:33], buf)
		}
	})

}

func TestWrite(t *testing.T) {
	file, err := os.Open("image.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	tFileName := filepath.Join(t.TempDir(), "oeinro.png")
	ofile, err := os.Create(tFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer ofile.Close()

	encode, err := NewPngEncoder(file, ofile, 0)
	if err != nil {
		t.Fatal(err)
	}

	encode.Write([]byte("hello world"))
}
