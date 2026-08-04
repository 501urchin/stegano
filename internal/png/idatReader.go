package png

import (
	"errors"
	"io"
	"slices"

	pngerrors "github.com/501urchin/stegano/v2/pkg/errors"
	"github.com/501urchin/stegano/v2/pkg/types"
)

type idatReader struct {
	imageChunks      []PngChunk
	srcReader        io.ReadSeeker
	currentChunk     int
	currentChunkSize uint32
	chunkOffset      uint32
}

func NewIdatReader(src io.ReadSeeker, chunks []PngChunk) (reader *idatReader, err error) {
	if len(chunks) <= 0 {
		return
	}

	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	if src == nil {
		return nil, pngerrors.ErrSourceIsNil
	}

	return &idatReader{
		imageChunks:      chunks,
		srcReader:        src,
		currentChunk:     0,
		currentChunkSize: 0,
		chunkOffset:      0,
	}, nil

}

func (r *idatReader) Close() (err error) {
	return nil
}

func (r *idatReader) prepareNextIdat() (err error) {
	for i := r.currentChunk + 1; i < len(r.imageChunks); i++ {
		if slices.Equal(r.imageChunks[i].Type, types.IEND) {
			return io.EOF
		}

		if !slices.Equal(r.imageChunks[i].Type, types.IDAT) {
			continue
		}

		r.currentChunk = i
		r.currentChunkSize = r.imageChunks[i].Length
		r.chunkOffset = 0

		_, err = r.srcReader.Seek(int64(r.imageChunks[i].DataStartIdx), io.SeekStart)
		if err != nil {
			return errors.Join(pngerrors.ErrFailedToSeek, err)
		}

		break
	}

	return nil
}

func (r *idatReader) Read(b []byte) (n int, err error) {
	if r.chunkOffset == 0 || r.currentChunkSize-1 <= r.chunkOffset {
		err = r.prepareNextIdat()
		if err != nil {
			return
		}
	}

	availableBytesInChunk := int(r.currentChunkSize - r.chunkOffset)

	if dstLen := len(b); dstLen <= availableBytesInChunk {
		r.chunkOffset += uint32(dstLen)
		return r.srcReader.Read(b)
	}

	firstPass, err := r.srcReader.Read(b[:availableBytesInChunk])
	if err != nil {
		return
	}
	r.chunkOffset += uint32(firstPass)

	err = r.prepareNextIdat()
	if err != nil {
		return
	}

	secondPass, err := r.srcReader.Read(b[availableBytesInChunk:])
	if err != nil {
		return
	}
	r.chunkOffset += uint32(secondPass)

	return
}
