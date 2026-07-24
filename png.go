package stegano

import (
	"fmt"
	"io"
	"slices"

	"github.com/501urchin/stegano/v2/internal/png"
	"github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type PngEncoder struct {
	src io.ReadSeeker
	dst io.WriteSeeker

	chunks []png.PngChunk

	embeddingDepth types.BitIndex
	pngInfo        png.PngInfo

	embeddingCapacity int64 // bytes
	remainigCapacity  int64 // bytes

	currentChunk int
	// TODO: add a reuseable buffer where we can store the idat line chunk
}

func (e *PngEncoder) Capacity() int64 {
	return e.embeddingCapacity
}

func (e *PngEncoder) Remaining() int64 {
	return e.remainigCapacity
}

func NewPngEncoder(src io.ReadSeeker, dst io.WriteSeeker, bitDepth types.BitIndex) (enc *PngEncoder, err error) {
	if bitDepth > 7 {
		return nil, errors.ErrInvalidBitDepth
	}

	if src == nil {
		return nil, errors.ErrSrcIsNil
	}

	if dst == nil {
		return nil, errors.ErrDstIsNil
	}

	enc = &PngEncoder{
		embeddingDepth: bitDepth,
		src:            src,
		dst:            dst,
	}

	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	enc.chunks, err = png.DecodeChunks(src)
	if err != nil {
		return
	}
	if len(enc.chunks) == 0 {
		return nil, errors.ErrMissingIHDR
	}

	if !slices.Equal(enc.chunks[0].Type, types.IHDR) {
		return nil, errors.ErrMissingIHDR
	}

	pngInfo, err := png.ParseIHDRChunk(src, enc.chunks[0])
	if err != nil {
		return
	}
	enc.pngInfo = pngInfo

	if uint8(bitDepth) > pngInfo.BitDepth-1 {
		return nil, errors.ErrBitDepthTooHigh
	}

	_, err = dst.Write(png.PngSignature)
	if err != nil {
		return
	}

	enc.currentChunk = -1
	for i, c := range enc.chunks {
		if slices.Equal(c.Type, types.IDAT) {
			if enc.currentChunk == -1 {
				enc.currentChunk = i
			}

			enc.embeddingCapacity += int64(c.Length / 8)
			continue
		}

		if enc.currentChunk == -1 {
			err = png.WriteChunk(src, dst, enc.chunks[i])
			if err != nil {
				return nil, errors.ErrFailedToWriteChunk
			}
		}
	}

	enc.remainigCapacity = enc.embeddingCapacity

	return
}

func (e *PngEncoder) Flush() (err error) {
	for i := e.currentChunk; i < len(e.chunks); i++ {
		err = png.WriteChunk(e.src, e.dst, e.chunks[i])
		if err != nil {
			return errors.ErrFailedToWriteChunk
		}
	}

	return
}

func (e *PngEncoder) Write(p []byte) (n int, err error) {
	if len(p) > int(e.remainigCapacity) {
		return 0, fmt.Errorf("no space left")
	}

	return
}


// IDEA: to reduce detectability we can pad the data with zero bytes causing the data to be spread apart
// this carries a few nuances
// since we dont know how much data the caller want to embed it will be hard to estimate how much padding to add. what if they want to embed more and they cant because we added to much padding
// second nuance is that we need a reliable way to add and deocde the padding. we need to create a way to introduce padding to where the bits are spread apart and that we can reverse back into a full piece of data