// Package png contains utilities for decoding png images
package png

import (
	"hash"
	"hash/crc32"
	"io"
)

type pngDecoder struct {
	crc        hash.Hash32
	imageInfo  ihdrData
	zlibReader io.ReadCloser
}

func NewPNGDecoder(src io.ReadSeeker) (d *pngDecoder, err error) {
	chunks, err := DecodeChunks(src)
	if err != nil {
		return
	}

	idatReader, err := NewIdatReader(src, chunks)
	if err != nil {
		return
	}

	ihdrData, err := ParseIHDRChunk(src, chunks[0])
	if err != nil {
		return
	}

	return &pngDecoder{
		crc:        crc32.NewIEEE(),
		imageInfo:  ihdrData,
		zlibReader: idatReader,
	}, nil
}
