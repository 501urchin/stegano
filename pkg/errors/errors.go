package errors

import "errors"

var (
	ErrInvalidBitDepth           = errors.New("bit depth must be between 0 and 7")
	ErrInvalidBitIndex           = errors.New("invalid bit index")
	ErrSrcIsNil                  = errors.New("source reader cannot be nil")
	ErrDstIsNil                  = errors.New("destination writer cannot be nil")
	ErrMissingIHDR               = errors.New("PNG must start with IHDR chunk")
	ErrBitDepthTooHigh           = errors.New("embedding bit depth exceeds PNG bit depth")
	ErrFailedToWriteChunk        = errors.New("failed to write PNG chunk")
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
	ErrSourceIsNil               = ErrSrcIsNil
	ErrInvalidChunk              = errors.New("chunk is invalid: length is zero")
	ErrFailedToReadChunkData     = errors.New("failed to read chunk data from source")
	ErrorFailedToReadNBytes      = errors.New("failed to read the exact amount of bytes from stream")
)
