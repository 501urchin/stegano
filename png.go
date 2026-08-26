package stegano

import (
	"compress/zlib"
	"io"
	"slices"

	"github.com/501urchin/stegano/v2/internal/bits"
	"github.com/501urchin/stegano/v2/internal/png"
	"github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type pngEncoder struct {
	imageHeader    png.IHDRData
	zlibReader     io.ReadCloser
	dstWriter      io.Writer
	rowBytes       int
	previousRow    []byte
	currentRow     []byte
	cursor         int
	bytesPerSample int
	bitDepth       int
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

	zlibReader, err := zlib.NewReader(idatReader)
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

			continue
		}

		if currentChunk == -1 {
			err = png.WriteChunk(src, dst, chunks[i])
			if err != nil {
				return nil, errors.ErrFailedToWriteChunk
			}
		}
	}

	bytesPerRow, err := png.BytePerRow(ihdrData.Width, ihdrData.Interlace, ihdrData.ColorType, ihdrData.BitDepth)
	if err != nil {
		return
	}

	enc := &pngEncoder{
		imageHeader:    ihdrData,
		zlibReader:     zlibReader,
		rowBytes:       bytesPerRow,
		previousRow:    make([]byte, bytesPerRow),
		currentRow:     make([]byte, bytesPerRow),
		dstWriter:      dst,
		bitDepth:       int(bitDepth),
		bytesPerSample: 1,
		cursor:         0,
	}

	if ihdrData.BitDepth == 16 {
		enc.bytesPerSample = 2
		enc.cursor = 1
	}

	return enc, nil
}

// read: zlib.Read -> interlace -> data
// write: data -> interlace -> zlib.Write

func (d *pngEncoder) Write(b []byte) (err error) {
	// BUG: problems here, idk what is causing it but it looks like the zlib reader returns incomplete data and the interlce reversal is not working
	d.zlibReader.Read(d.currentRow)
	bpp := png.BitsPerPixel(d.imageHeader.ColorType, int(d.imageHeader.BitDepth))
	d.currentRow = png.ReverseNoneFilter(d.currentRow, d.previousRow, bpp)

	for _, v := range b {
		for i := 0; i < 8; i++ {
			dif := bits.At(v, i)
			if dif != bits.At(d.currentRow[d.cursor], d.bitDepth) {
				d.currentRow[d.cursor] = bits.Flip(d.currentRow[d.cursor], d.bitDepth)
			}
			d.cursor += d.bytesPerSample

		}
	}

	return
}

// 101
// 001
