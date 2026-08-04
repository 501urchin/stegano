package stegano

import (
	"hash"
	"hash/crc32"
	"io"

	"github.com/501urchin/stegano/v2/internal/png"
)

type pngDecoder struct {
	crc        hash.Hash32
	imageInfo  png.IHDRData
	zlibReader io.ReadCloser
}

func NewPNGDecoder(src io.ReadSeeker) (d *pngDecoder, err error) {
	chunks, err := png.DecodeChunks(src)
	if err != nil {
		return
	}

	idatReader, err := png.NewIdatReader(src, chunks)
	if err != nil {
		return
	}

	ihdrData, err := png.ParseIHDRChunk(src, chunks[0])
	if err != nil {
		return
	}

	return &pngDecoder{
		crc:        crc32.NewIEEE(),
		imageInfo:  ihdrData,
		zlibReader: idatReader,
	}, nil
}

func (d *pngDecoder) Embed(b []byte) (err error) {
	d.zlibReader.Read()
	return
}
