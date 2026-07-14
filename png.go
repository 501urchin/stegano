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
func NewPngEncoder(src io.ReadSeeker, dst io.WriteSeeker, bitDepth ...types.BitIndex) (enc *PngEncoder, err error)

func (e *PngEncoder) Close() (err error)
func (e *PngEncoder) Write(data []byte)
