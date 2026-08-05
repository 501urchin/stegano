package png

import (
	"fmt"
	"math"
)

func BytePerRow(width uint32, interlace, color, bitdepth uint8) (int, error) {
	if interlace != 0 {
		return 0, fmt.Errorf("no support for adam7 interlace")
	}

	var channels int

	switch color {
	case 0:
		channels = 1
	case 2:
		channels = 3
	case 3:
		channels = 1
	case 4:
		channels = 2
	case 6:
		channels = 4
	default:
		return 0, fmt.Errorf("invalid color")
	}

	bits := int(width) * channels * int(bitdepth)

	// round up bits to bytes
	bytes := int(math.Ceil(float64(bits) / 8))

	// filter byte
	return bytes + 1, nil
}

func ReverseNoneFilter(currentRow, previousRow []byte, bpp int) []byte {
	if len(currentRow) == 0 {
		return nil
	}

	filter := currentRow[0]
	out := make([]byte, len(currentRow))
	out[0] = filter

	switch filter {

	case 0:
		copy(out, currentRow)
	case 1:
		for i := 1; i < len(currentRow); i++ {
			left := byte(0)

			if i > bpp {
				left = out[i-bpp]
			}

			out[i] = currentRow[i] + left
		}

	case 2:
		for i := 1; i < len(currentRow); i++ {
			up := byte(0)

			if len(previousRow) > 0 {
				up = previousRow[i]
			}

			out[i] = currentRow[i] + up
		}

	case 3:
		for i := 1; i < len(currentRow); i++ {
			left := byte(0)
			up := byte(0)

			if i > bpp {
				left = out[i-bpp]  
			}

			if len(previousRow) > 0 {
				up = previousRow[i]
			}

			avg := byte((int(left) + int(up)) / 2)

			out[i] = currentRow[i] + avg
		}

	case 4:
		for i := 1; i < len(currentRow); i++ {
			left := byte(0)
			up := byte(0)
			upLeft := byte(0)

			if i > bpp {
				left = out[i-bpp]
			}

			if len(previousRow) > 0 {
				up = previousRow[i]

				if i > bpp {
					upLeft = previousRow[i-bpp]
				}
			}

			out[i] = currentRow[i] + paeth(left, up, upLeft)
		}

	default:
		panic("invalid png filter")
	}

	return out
}

func paeth(a, b, c byte) byte {
	p := int(a) + int(b) - int(c)

	pa := math.Abs(float64(p - int(a)))
	pb := math.Abs(float64(p - int(b)))
	pc := math.Abs(float64(p - int(c)))

	if pa <= pb && pa <= pc {
		return a
	}

	if pb <= pc {
		return b
	}

	return c
}
