package png

import "encoding/binary"

// ParseIHDRChunk takes in the data line of the IHDR chunk
func ParseIHDRChunk(data []byte) (PngInfo, error) {
	if len(data) != 13 {
		return PngInfo{}, ErrNotIHDR
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
		return PngInfo{}, ErrInvalidPNGBitDepth
	}
	
	switch info.ColorType {
	case 0, 2, 3, 4, 6:
		break
	default:
		return PngInfo{}, ErrInvalidPNGColorType
	}

	if info.Compression != 0 {
		return PngInfo{}, ErrInvalidPNGCompression
	}

	if info.Filter != 0 {
		return PngInfo{}, ErrInvalidPNGFilter
	}

	if info.Interlace > 1 {
		return PngInfo{}, ErrInvalidPNGInterlace
	}

	return info, nil
}
