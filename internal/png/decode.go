// Package png contains utilities for decoding png images
package png

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"slices"
)

func DecodePNG(file io.ReadSeeker) (chunks []PngChunk, err error) {
	if file == nil {
		return nil, ErrNotPNG
	}

	var buf = make([]byte, 8)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return nil, errors.Join(ErrFailedToReadSignature, err)
	}

	if !slices.Equal(buf, PngSignature) {
		return nil, ErrNotPNG
	}

	for {
		_, err = io.ReadFull(file, buf[:8])
		if err != nil && !errors.Is(err, io.EOF) {
			return chunks, errors.Join(ErrFailedToReadTypeAndLength, err)
		}

		c := PngChunk{
			Length: binary.BigEndian.Uint32(buf[:4]),
			Type:   append([]byte(nil), buf[4:8]...),
		}

		dataStartIdx, intErr := file.Seek(0, io.SeekCurrent)
		if intErr != nil {
			return chunks, errors.Join(ErrFailedToSeek, err)
		}

		c.DataStartIdx = int(dataStartIdx)

		_, err = file.Seek(int64(c.Length), io.SeekCurrent)
		if err != nil {
			return chunks, errors.Join(ErrFailedToSeek, err)
		}
		
		_, err = io.ReadFull(file, buf[:4])
		if err != nil {
			return chunks, errors.Join(ErrFailedToReadCRC, err)
		}

		c.CRC = binary.BigEndian.Uint32(buf[:4])

		chunks = append(chunks, c)
		if string(c.Type) == "IEND" {
			break
		}
	}

	var biggestChunkLen int
	for _, c := range chunks {
		if c.Length > uint32(biggestChunkLen) {
			biggestChunkLen = int(c.Length)
		}
	}

	dataBuf := make([]byte, biggestChunkLen)
	h := crc32.NewIEEE()
	for _, c := range chunks {
		_, err = file.Seek(int64(c.DataStartIdx), io.SeekStart)
		if err != nil {
			return
		}

		_, err = io.ReadFull(file, dataBuf[:c.Length])
		if err != nil {
			return
		}

		h.Reset()
		_, err = h.Write(c.Type)
		if err != nil {
			return
		}

		_, err = h.Write(dataBuf[:c.Length])
		if err != nil {
			return
		}

		if h.Sum32() != c.CRC {
			return nil, ErrCRCMismatch
		}
	}

	return
}
