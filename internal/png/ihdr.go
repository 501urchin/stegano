package png

import (
	"encoding/binary"
	"io"

	pngerrors "github.com/501urchin/stegano/v2/pkg/errors"
)

// ParseIHDRChunk takes in the data line of the IHDR chunk
func ParseIHDRChunk(src io.ReadSeeker, c PngChunk) (PngInfo, error) {
	data, err := GetChunkData(src, c)
	if err != nil {
		return PngInfo{}, err
	}

	info := PngInfo{
		Width:       binary.BigEndian.Uint32(data[0:4]),
		Height:      binary.BigEndian.Uint32(data[4:8]),
		BitDepth:    data[8],
		ColorType:   data[9],
		Compression: data[10],
		Filter:      data[11],
		Interlace:   data[12],
	}

	switch info.BitDepth {
	case 1, 2, 4, 8, 16:
		break
	default:
		return PngInfo{}, pngerrors.ErrInvalidPNGBitDepth
	}

	switch info.ColorType {
	case 0, 2, 3, 4, 6:
		break
	default:
		return PngInfo{}, pngerrors.ErrInvalidPNGColorType
	}

	if info.Compression != 0 {
		return PngInfo{}, pngerrors.ErrInvalidPNGCompression
	}

	if info.Filter != 0 {
		return PngInfo{}, pngerrors.ErrInvalidPNGFilter
	}

	if info.Interlace > 1 {
		return PngInfo{}, pngerrors.ErrInvalidPNGInterlace
	}

	return info, nil
}
