package png

import (
	"encoding/binary"
	"io"

	pngerrors "github.com/501urchin/stegano/v2/pkg/errors"
)

// ParseIHDRChunk takes in the data line of the IHDR chunk
func ParseIHDRChunk(src io.ReadSeeker, c PngChunk) (IHDRData, error) {
	data, err := GetChunkData(src, c)
	if err != nil {
		return IHDRData{}, err
	}

	info := IHDRData{
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
		return IHDRData{}, pngerrors.ErrInvalidPNGBitDepth
	}

	switch info.ColorType {
	case 0, 2, 3, 4, 6:
		break
	default:
		return IHDRData{}, pngerrors.ErrInvalidPNGColorType
	}

	if info.Compression != 0 {
		return IHDRData{}, pngerrors.ErrInvalidPNGCompression
	}

	if info.Filter != 0 {
		return IHDRData{}, pngerrors.ErrInvalidPNGFilter
	}

	if info.Interlace > 1 {
		return IHDRData{}, pngerrors.ErrInvalidPNGInterlace
	}

	return info, nil
}
