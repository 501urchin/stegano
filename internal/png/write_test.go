package png

import (
	"bytes"
	"crypto/rand"
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
		t.Errorf("result didnt match expecting: \n%v\n%v\n", expecting, out.Bytes())
	}

}

func BenchmarkWriteChunk(b *testing.B) {
	data := make([]byte, 1 << 20)
	_, err := rand.Read(data)
	if err != nil {
		b.Fatal(err)
	}
	source := bytes.NewReader([]byte{1, 2})
	c := PngChunk{
		Length:       1 << 20,
		Type:         []byte("IDAT"),
		DataStartIdx: 0,
	}
	out := bytes.NewBuffer(nil)

	b.ResetTimer()
	for b.Loop() {
		_ = WriteChunk(source, out, c)
	}
}
