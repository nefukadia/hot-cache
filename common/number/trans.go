package number

type Integer interface {
	~uint | ~uintptr | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

// LowByteToOther get low byte to targetType
func LowByteToOther[sourceType Integer, targetType Integer](x sourceType) targetType {
	return targetType(byte(x))
}
