// Package png contains utilities for handling png images
package png

type PngChunk struct {
	Length       uint32
	Type         []byte
	DataStartIdx int // start index of the data line relative to the start of the file
	CRC          uint32
}

type IHDRData struct {
	Width       uint32
	Height      uint32
	BitDepth    uint8 // 1, 2, 4, 8, 16
	ColorType   uint8 // 0, 2, 3, 4, 6
	Compression uint8
	Filter      uint8
	Interlace   uint8
}

var (
	PngSignature = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
)
