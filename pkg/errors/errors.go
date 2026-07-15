package errors

import "errors"

var (
	ErrInvalidBitDepth = errors.New("bit depth must be between 0 and 7")
	ErrSrcIsNil        = errors.New("source reader cannot be nil")
	ErrDstIsNil        = errors.New("destination writer cannot be nil")
	ErrMissingIHDR     = errors.New("PNG must start with IHDR chunk")
	ErrBitDepthTooHigh = errors.New("embedding bit depth exceeds PNG bit depth")
)
