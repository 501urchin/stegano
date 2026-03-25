package stegano

import (
	"errors"
	"fmt"
	"io"

	"github.com/501urchin/stegano/v2/internal/png"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type PngEncoder struct {
	carrier             io.ReadSeeker
	carrierChunks       []png.PngChunk
	out                 io.Writer
	bitdepth            types.BitIndex
	operatingChunkIndex int
	capacity            int
	written             int
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

	chunks, err := png.DecodePNG(carrier)
	if err != nil {
		return
	}
	enc.carrierChunks = chunks

	for _, c := range chunks {
		if string(c.Type) == "IDAT" {
			enc.capacity += int(c.Length) / 8
		}
	}

	written, err := out.Write(png.PngSignature)
	if err != nil || written != 8 {
		return nil, fmt.Errorf("failed to write png header to out: %v", err)
	}

	err = png.WriteChunk(chunks[0], carrier, out)
	if err != nil {
		return
	}

	enc.operatingChunkIndex++

	return enc, nil
}
func (e *PngEncoder) Capacity() int64    { return int64(e.capacity) }
func (e *PngEncoder) Remaining() int64   { return int64(e.capacity) - int64(e.written) }
func (e *PngEncoder) Close() (err error) { return }

func (e *PngEncoder) Write(data []byte) (written int, err error) { return }
