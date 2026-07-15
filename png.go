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
	embeddingCapacity int
}

func (e *PngEncoder) Capacity() int64

func (e *PngEncoder) Remaining() int64
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
	}

	enc.chunks, err = png.ParsePNG(src)
	if err != nil {
		return
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

	

	for _, c := enc.chunks {
		if c.
	}

	// TODO: get steganography capacity. each rgb pixel can hold 1 bit
	// TODO: write png, ihdr and any other header before the first idat to the dst
	return
}

func (e *PngEncoder) Close() (err error)

// TODO: remember png idat capacity is influenced by bit depth in the ihdr header. some pngs can be 16 bit per channel or 1 bit per channel. factor this in when processing the idat lines and encoding the bits
func (e *PngEncoder) Write(data []byte)
