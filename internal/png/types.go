package png

import "errors"

type PngChunk struct {
	Length       uint32
	Type         []byte
	DataStartIdx int // start index of the data line relative to the start of the file
	CRC          uint32
}

type PngInfo struct {
	Width       uint32
	Height      uint32
	BitDepth    uint8 // 1, 2, 4, 8, 16
	ColorType   uint8 // 0, 2, 3, 4, 6
	Compression uint8
	Filter      uint8
	Interlace   uint8
}

var (
	PngSignature                 = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	ErrorFailedToReadNBytes      = errors.New("failed to read the exact amount of bytes from stream")
	ErrCRCMismatch               = errors.New("crc doesnt match data and type")
	ErrNotPNG                    = errors.New("file does not have a valid png signature")
	ErrFailedToReadTypeAndLength = errors.New("failed to read png chunk type and length")
	ErrFailedToSeek              = errors.New("failed to seek png at index")
	ErrFailedToReadCRC           = errors.New("failed to read png chunk crc")
	ErrFailedToReadSignature     = errors.New("failed to read 8 byte png signature")
	ErrNotIHDR                   = errors.New("data is not valid for IHDR")
	ErrInvalidPNGBitDepth        = errors.New("IHDR header contains an invalid bit depth")
	ErrInvalidPNGColorType       = errors.New("IHDR header contains an invalid color type")
	ErrInvalidPNGCompression     = errors.New("IHDR header contains an invalid compression method")
	ErrInvalidPNGFilter          = errors.New("IHDR header contains an invalid filter method")
	ErrInvalidPNGInterlace       = errors.New("IHDR header contains an invalid interlace method")
	ErrSourceIsNil               = errors.New("source reader cannot be nil")
	ErrInvalidChunk              = errors.New("chunk is invalid: length is zero")
	ErrFailedToReadChunkData     = errors.New("failed to read chunk data from source")
)
