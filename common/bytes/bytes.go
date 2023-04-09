package bytes

import (
	"errors"
	"math"
)

func VLEMarshal(data []byte) []byte {
	if len(data) == 0 {
		return make([]byte, 0)
	}

	ret := make([]byte, int(math.Ceil(float64(len(data)*8)/7.0)))

	for i := 0; i < len(ret); i++ {
		var startByte = i * 7 / 8
		var offsetBit = (i * 7) % 8
		var markHigh, markLow byte
		switch offsetBit {
		case 0:
			markHigh = 1<<8 - 2
		default:
			markHigh = 1<<(8-offsetBit) - 1
			if startByte+1 < len(data) {
				markLow = byte(0xff - (1<<(9-offsetBit) - 1))
			}
		}

		if markLow > 0 {
			ret[i] = 1<<7 + (data[startByte]&markHigh)<<(offsetBit-1)
			if startByte+1 < len(data) {
				ret[i] += (data[startByte+1] & markLow) >> (9 - offsetBit)
			}
		} else {
			tmp := data[startByte] & markHigh

			// markHigh = 11111110
			if markHigh&byte(0x80) != 0 {
				tmp = tmp >> 1
			} else {
				for markHigh&byte(0x40) == 0 {
					markHigh = markHigh << 1
					tmp = tmp << 1
				}
			}
			ret[i] = 1<<7 + tmp
		}
	}
	ret[len(ret)-1] = ret[len(ret)-1] & 0x7f

	return ret
}

func VLEUnmarshal(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return make([]byte, 0), nil
	}
	if len(data) <= 1 || data[len(data)-1]&byte(0x80) != 0 || ((len(data)-1)*7)%8 == 0 {
		return nil, errors.New("bad data")
	}

	for i, b := range data {
		if i+1 != len(data) && (b&0x80) == 0 {
			return nil, errors.New("bad data in forward byte")
		}
	}

	ret := make([]byte, int(math.Ceil(float64((len(data)-1)*7)/8.0)))

	for i := range ret {
		var startByte = int(math.Floor(float64(i*8) / 7.0))
		var offsetBit = (i*8)%7 + 1
		var markHigh, markLow byte
		markHigh = 0xff >> offsetBit
		markLow = 0xff << (7 - offsetBit) & 0x7f
		ret[i] = (data[startByte]&markHigh)<<offsetBit + (data[startByte+1]&markLow)>>(7-offsetBit)
	}

	return ret, nil
}
