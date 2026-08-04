// Package png contains utilities for decoding png images
package png

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"slices"

	pngerrors "github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

func NewPNGDecoder() pngDecoder {
	return pngDecoder{
		check:     crc32.NewIEEE(),
		imageInfo: ihdrData{},
	}
}

func validateChunk(chunk PngChunk, file io.ReadSeeker) (err error) {
	h := crc32.NewIEEE()
	_, err = file.Seek(int64(chunk.DataStartIdx), io.SeekStart)
	if err != nil {
		return
	}

	h.Reset()
	_, err = h.Write(chunk.Type)
	if err != nil {
		return
	}

	_, err = io.CopyN(h, file, int64(chunk.Length))
	if err != nil {
		return
	}

	if h.Sum32() != chunk.CRC {
		return pngerrors.ErrCRCMismatch
	}

	return nil
}

func DecodeChunks(file io.ReadSeeker) (chunks []PngChunk, err error) {
	if file == nil {
		return nil, pngerrors.ErrNotPNG
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	var buf = make([]byte, 8)

	// read png signature
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return nil, errors.Join(pngerrors.ErrFailedToReadSignature, err)
	}

	if !slices.Equal(buf, PngSignature) {
		return nil, pngerrors.ErrNotPNG
	}

	var biggestChunkLen int
	for {
		// get len and type
		_, err = io.ReadFull(file, buf[:8])
		if err != nil && !errors.Is(err, io.EOF) {
			return chunks, errors.Join(pngerrors.ErrFailedToReadTypeAndLength, err)
		}

		if errors.Is(err, io.EOF) {
			if len(chunks) == 0 {
				return nil, pngerrors.ErrMissingIHDR
			}

			break
		}

		chunkLength := binary.BigEndian.Uint32(buf[:4])
		if chunkLength > uint32(biggestChunkLen) {
			biggestChunkLen = int(chunkLength)
		}

		c := PngChunk{
			Length: chunkLength,
			Type:   append([]byte(nil), buf[4:8]...),
		}

		// get start index of idat line for this chunk
		dataStartIdx, intErr := file.Seek(0, io.SeekCurrent)
		if intErr != nil {
			return chunks, errors.Join(pngerrors.ErrFailedToSeek, err)
		}
		c.DataStartIdx = int(dataStartIdx)

		_, err = file.Seek(int64(c.Length), io.SeekCurrent)
		if err != nil {
			return chunks, errors.Join(pngerrors.ErrFailedToSeek, err)
		}

		// get crc
		_, err = io.ReadFull(file, buf[:4])
		if err != nil {
			return chunks, errors.Join(pngerrors.ErrFailedToReadCRC, err)
		}

		c.CRC = binary.BigEndian.Uint32(buf[:4])

		// TODO: figure out if there is a way to estimate how many chunks a png has so we can redouce allocs
		chunks = append(chunks, c)
		if slices.Equal(c.Type, types.IEND) {
			break
		}
	}


	for _, c := range chunks {
		err = validateChunk(c, file)
		if err != nil {
			return nil, err
		}
	}

	return
}

func GetChunkData(src io.ReadSeeker, c PngChunk) (data []byte, err error) {
	if src == nil {
		return nil, pngerrors.ErrSourceIsNil
	}

	if c.Length == 0 {
		return nil, pngerrors.ErrInvalidChunk
	}

	_, err = src.Seek(int64(c.DataStartIdx), io.SeekStart)
	if err != nil {
		return nil, errors.Join(pngerrors.ErrFailedToSeek, err)
	}

	data = make([]byte, c.Length)

	_, err = io.ReadFull(src, data)
	if err != nil {
		return nil, errors.Join(pngerrors.ErrFailedToReadChunkData, err)
	}

	return data, nil
}
