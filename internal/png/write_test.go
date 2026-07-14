package png

import (
	"bytes"
	"slices"
	"testing"
)

func TestWriteChunk(t *testing.T) {
	source := bytes.NewReader([]byte{1, 2})
	c := PngChunk{
		Length:       2,
		Type:         []byte("IDAT"),
		DataStartIdx: 0,
		CRC:          2,
	}

	out := bytes.NewBuffer(nil)

	err := WriteChunk(source, out, c)
	if err != nil {
		t.Fatal(err)
	}
	expecting := []byte{0, 0, 0, 2, 73, 68, 65, 84, 1, 2, 139, 238, 237, 215}

	if !slices.Equal(out.Bytes(), expecting) {
		t.Error("result didnt match expecting")
	}

}

func BenchmarkWriteChunk(b *testing.B) {
	source := bytes.NewReader([]byte{1, 2})
	c := PngChunk{
		Length:       2,
		Type:         []byte("IDAT"),
		DataStartIdx: 0,
		CRC:          2,
	}
	out := bytes.NewBuffer(nil)

	for b.Loop() {
		_ = WriteChunk(source, out, c)
	}
}
