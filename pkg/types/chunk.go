package types

type ChunkType []byte

var (
	IHDR ChunkType = []byte("IHDR")
	IDAT ChunkType = []byte("IDAT")
)
