package png

import (
	"encoding/binary"
	"hash/crc32"
	"io"
)

func WriteChunk(source io.ReadSeeker, out io.Writer, chunks ...PngChunk) (err error) {
	if len(chunks) == 0 {
		return
	}

	var biggestChunkLen int
	for _, c := range chunks {
		if c.Length > uint32(biggestChunkLen) {
			biggestChunkLen = int(c.Length)
		}
	}

	buf := make([]byte, max(biggestChunkLen, 4))

	for _, chunk := range chunks {
		binary.BigEndian.PutUint32(buf, chunk.Length)
		_, err = out.Write(buf[:4])
		if err != nil {
			return
		}

		_, err = out.Write(chunk.Type)
		if err != nil {
			return
		}

		_, err = source.Seek(int64(chunk.DataStartIdx), io.SeekStart)
		if err != nil {
			return
		}

		_, err = io.ReadFull(source, buf[:chunk.Length])
		if err != nil {
			return
		}

		_, err = out.Write(buf[:chunk.Length])
		if err != nil {
			return
		}

		h := crc32.NewIEEE()
		_, err = h.Write(chunk.Type)
		if err != nil {
			return
		}

		_, err = h.Write(buf[:chunk.Length])
		if err != nil {
			return
		}

		buf = buf[:0]
		buf = h.Sum(buf)

		_, err = out.Write(buf[:4])
		if err != nil {
			return
		}
	}

	return
}
