package internal

import (
	"io"

	"github.com/501urchin/stegano/v2/pkg/types"
)

type PngDecoder struct {
	carrier          io.ReadSeeker
	bitdepth         types.BitIndex
	reachedEndOfData bool
}

func (d *PngDecoder) Read(data []byte) (read int, err error) { return }
func (e *PngDecoder) DataLength() int64
