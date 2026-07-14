package stegano

import (
	"errors"
	"io"

	"github.com/501urchin/stegano/v2/internal/png"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type PngEncoder struct {
	src    io.ReadSeeker
	chunks []png.PngChunk
	dst    io.WriteSeeker

	bitDepth types.BitIndex
}

var (
	ErrInvalidBitDepth = errors.New("invalid bit depth")
)

func (e *PngEncoder) Capacity() int64

func (e *PngEncoder) Remaining() int64
func NewPngEncoder(src io.ReadSeeker, dst io.WriteSeeker, bitDepth ...types.BitIndex) (enc *PngEncoder, err error) {
	// TODO: validate bit depth
	// TODO: get steganography capacity. each rgb pixel can hold 1 bit
	// TODO: write png, ihdr and any other header before the first idat to the dst
	return
}

func (e *PngEncoder) Close() (err error)

// TODO: remember png idat capacity is influenced by bit depth in the ihdr header. some pngs can be 16 bit per channel or 1 bit per channel. factor this in when processing the idat lines and encoding the bits
func (e *PngEncoder) Write(data []byte)
