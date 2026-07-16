package stegano

import (
	"os"
	"testing"

	"github.com/501urchin/stegano/v2/pkg/types"
)

// TestCase for PngEncode
// test for invalid params
// test for invalid png (headder, ihdr)
// test for idat line in different spot eg chunk index 10 instead of 0

// TestCase for Flush
// test if dst image fully matches src
// test flush at middle of image

func TestPngEncode(t *testing.T) {
	file, err := os.Open("image.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	nf, err := os.Create("outfile.png")
	if err != nil {
		t.Fatal(err)
	}
	defer nf.Close()

	enc, err := NewPngEncoder(file, nf, types.BitZero)
	if err != nil {
		panic(err)
	}

	err = enc.Flush()
	if err != nil {
		t.Fatal(err)
	}
}
