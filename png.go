package stegano

import (
	"errors"
	"io"

	"github.com/501urchin/stegano/v2/pkg/types"
)

type PngEncoder struct {
	carrier  io.ReadSeeker
	out      io.Writer
	bitdepth types.BitIndex
}

var (
	ErrInvalidBitDepth = errors.New("invalid bit depth")
)

func NewPngEncoder(carrier io.ReadSeeker, out io.Writer, bitdepth ...types.BitIndex) (enc *PngEncoder, err error) {
	providedBitDepth := len(bitdepth) != 0
	if providedBitDepth && bitdepth[0] > 7 {
		return nil, ErrInvalidBitDepth
	}

	enc = &PngEncoder{
		carrier: carrier,
		out:     out,
	}

	if providedBitDepth {
		enc.bitdepth = bitdepth[0]
	} else {
		enc.bitdepth = types.BitOne
	}

	return enc, nil
}

func (e *PngEncoder) Write(data []byte) (written int, err error) { return }
func (e *PngEncoder) Capacity() int64                            { return 0 }
func (e *PngEncoder) Remaining() int64                           { return 0 }
func (e *PngEncoder) Close() (err error)                         { return }
