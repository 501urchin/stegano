package png

import "errors"

type pngChunk struct {
	Length       uint32
	Type         []byte
	DataStartIdx int // start index of the data line relative to the start of the file
	CRC          uint32
}

var (
	pngSignature                 = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	ErrorFailedToReadNBytes      = errors.New("failed to read the exact amount of bytes from stream")
	ErrCRCMismatch               = errors.New("crc doesnt match data and type")
	ErrNotPNG                    = errors.New("file does not have a valid png signature")
	ErrFailedToReadTypeAndLength = errors.New("failed to read png chunk type and length: ")
	ErrFailedToSeek              = errors.New("failed seek png at index: ")
	ErrFailedToReadCRC           = errors.New("failed to read png chunk crc: ")
	ErrFailedToReadSignature     = errors.New("failed to read 8 byte png signature: ")
)
