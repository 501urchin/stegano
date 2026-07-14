package png

import "encoding/binary"

// TODO: write test
func ParseIHDRChunk(data []byte) (d PngInfo, err error) {
	if len(data) != 13 {
		return d, ErrNotIHDR
	}

	d.Width = binary.BigEndian.Uint32(data[:4])
	d.Height = binary.BigEndian.Uint32(data[4:8])
	d.BitDepth = data[8]
	d.ColorType = data[9]
	d.Compression = data[10]
	d.Filter = data[11]
	d.Interlace = data[12]

	return
}
