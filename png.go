package stegano

import (
	"io"
	"slices"

	"github.com/501urchin/stegano/v2/internal/png"
	"github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type PngEncoder struct {
	src               io.ReadSeeker
	dst               io.WriteSeeker
	chunks            []png.PngChunk
	embeddingDepth    types.BitIndex
	pngInfo           png.PngInfo
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

	firstIDATChunkIndex := -1
	for i, c := range enc.chunks {
		if !slices.Equal(c.Type, types.IDAT) {
			continue
		}

		if firstIDATChunkIndex == -1 {
			firstIDATChunkIndex = i
		}

		enc.embeddingCapacity += int64(c.Length / 8)
	}
	enc.remainigCapacity = enc.embeddingCapacity

	_, err = dst.Write(png.PngSignature)
	if err != nil {
		return
	}

	for i := 0; i < firstIDATChunkIndex; i++ {
		err = png.WriteChunk(src, dst, enc.chunks[i])
		if err != nil {
			return nil, errors.ErrFailedToWriteChunk
		}
	}
	enc.currentChunk = firstIDATChunkIndex

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

// TODO: remember png idat capacity is influenced by bit depth in the ihdr header. some pngs can be 16 bit per channel or 1 bit per channel. factor this in when processing the idat lines and encoding the bits
// func (e *PngEncoder) Write(p []byte) (n int, err error) {
// 	return
// }
