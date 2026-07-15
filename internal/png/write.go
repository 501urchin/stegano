package png

import (
	"encoding/binary"
	"hash/crc32"
	"io"
)

func WriteChunk(source io.ReadSeeker, out io.Writer, chunk PngChunk) (err error) {
	_, err = out.Write(binary.BigEndian.AppendUint32(nil, chunk.Length))
	if err != nil {
		return
	}

	_, err = out.Write(chunk.Type)
	if err != nil {
		return
	}

	// copy data from src to dst
	_, err = source.Seek(int64(chunk.DataStartIdx), io.SeekStart)
	if err != nil {
		return
	}

	_, err = io.CopyN(out, source, int64(chunk.Length))
	if err != nil {
		return
	}

	// calc and write crc
	h := crc32.NewIEEE()
	_, err = h.Write(chunk.Type)
	if err != nil {
		return
	}

	_, err = source.Seek(int64(chunk.DataStartIdx), io.SeekStart)
	if err != nil {
		return
	}

	_, err = io.CopyN(h, source, int64(chunk.Length))
	if err != nil {
		return
	}

	_, err = out.Write(h.Sum(nil))
	if err != nil {
		return
	}

	return
}
