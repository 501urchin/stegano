package png

import (
	"encoding/binary"
	"hash/crc32"
	"io"
)

func WriteChunk(chunk PngChunk, source io.ReadSeeker, out io.Writer) (err error) {
	buf := make([]byte, chunk.Length)

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

	_, err = io.ReadAll(source)
	if err != nil {
		return
	}

	_, err = out.Write(buf)
	if err != nil {
		return
	}

	h := crc32.NewIEEE()
	_, err = h.Write(chunk.Type)
	if err != nil {
		return
	}

	_, err = h.Write(buf)
	if err != nil {
		return
	}

	_, err = out.Write(h.Sum(nil))
	if err != nil {
		return
	}

	return
}
