package png

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Dont use this function. its only for testing
func encode(newFileName, oldFileName string, chunks []pngChunk) (err error) {
	newFile, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	oldFile, err := os.Open(oldFileName)
	if err != nil {
		return
	}
	defer oldFile.Close()

	_, err = newFile.Write(pngSignature)
	if err != nil {
		return err
	}

	uint32Buf := make([]byte, 4)
	for i := range chunks {
		binary.BigEndian.PutUint32(uint32Buf, chunks[i].Length)
		_, err = newFile.Write(uint32Buf)
		if err != nil {
			return err
		}

		_, err = newFile.Write(chunks[i].Type)
		if err != nil {
			return err
		}

		_, err = oldFile.Seek(int64(chunks[i].DataStartIdx), io.SeekStart)
		if err != nil {
			return
		}

		if written, cErr := io.CopyN(newFile, oldFile, int64(chunks[i].Length)); cErr != nil || written != int64(chunks[i].Length) {
			return fmt.Errorf("failed to copy correct amount of bytes")
		}

		binary.BigEndian.PutUint32(uint32Buf, chunks[i].CRC)
		_, err = newFile.Write(uint32Buf)
		if err != nil {
			return err
		}
	}

	return nil
}
