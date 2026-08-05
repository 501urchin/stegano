package png

import (
	"fmt"
	"math"
)

func BytePerRow(width uint32, interlace, color, bitdepth uint8) (bytes int, err error) {
	if interlace == 0 {
		switch color {
		case 0, 3:
			bytes = 1
		case 2:
			bytes = 3
		case 4:
			bytes = 2
		case 6:
			bytes = 4
		default:
			return 0, fmt.Errorf("invalid color")
		}

		bytes *= int(bitdepth)
		bytes *= int(width)

		bytes = int(math.Ceil(float64(bytes) / float64(8)))

		bytes += 1
	} else {
		return 0, fmt.Errorf("no support for adams7 interlace")
	}

	return
}
