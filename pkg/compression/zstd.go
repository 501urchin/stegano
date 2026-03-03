package compression

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
)

func CompressZSTD(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	encoder, err := zstd.NewWriter(&buf)
	if err != nil {
		return nil, err
	}

	_, err = encoder.Write(data)
	if err != nil {
		return nil, err
	}

	encoder.Close()
	return buf.Bytes(), err
}

func DecompressZSTD(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	decoder, err := zstd.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	var out bytes.Buffer
	if _, err = io.Copy(&out, decoder); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
