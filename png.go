package stegano

import (
	"hash"
	"hash/crc32"
	"io"
	"slices"

	"github.com/501urchin/stegano/v2/internal/png"
	"github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type pngEncoder struct {
	crc        hash.Hash32
	imageInfo  png.IHDRData
	zlibReader io.ReadCloser
	dstWriter  io.Writer
}

func NewPngEncoder(src io.ReadSeeker, dst io.Writer, bitDepth types.BitIndex) (d *pngEncoder, err error) {
	if bitDepth > 7 {
		return nil, errors.ErrInvalidBitDepth
	}

	if src == nil {
		return nil, errors.ErrSrcIsNil
	}

	if dst == nil {
		return nil, errors.ErrDstIsNil
	}

	chunks, err := png.DecodeChunks(src)
	if err != nil {
		return
	}

	if len(chunks) == 0 {
		return nil, errors.ErrMissingIHDR
	}

	if !slices.Equal(chunks[0].Type, types.IHDR) {
		return nil, errors.ErrMissingIHDR
	}

	idatReader, err := png.NewIdatReader(src, chunks)
	if err != nil {
		return
	}

	ihdrData, err := png.ParseIHDRChunk(src, chunks[0])
	if err != nil {
		return
	}

	_, err = dst.Write(png.PngSignature)
	if err != nil {
		return
	}

	currentChunk := -1
	for i, c := range chunks {
		if slices.Equal(c.Type, types.IDAT) {
			if currentChunk == -1 {
				currentChunk = i
			}

			// embeddingCapacity += int64(c.Length / 8)
			continue
		}

		if currentChunk == -1 {
			err = png.WriteChunk(src, dst, chunks[i])
			if err != nil {
				return nil, errors.ErrFailedToWriteChunk
			}
		}
	}

	return &pngEncoder{
		crc:        crc32.NewIEEE(),
		imageInfo:  ihdrData,
		zlibReader: idatReader,
		dstWriter:  dst,
	}, nil
}

func (d *pngEncoder) Write(b []byte) (err error) {
	
	// d.zlibReader.Read()
	return
}
