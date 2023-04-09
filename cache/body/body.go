package body

const LenByte = 8
const CfgIntegerByte = 8

type Body interface {
	GetLen() int64
	GetCfg() map[byte]byte
	GetData() map[byte]any
	Valid() bool
	ToBytes() ([]byte, error)
	Setup([]byte) error
}

type BaseBody struct {
	Len  int64
	Cfg  map[byte]byte
	Data map[byte]any
}

func (bb *BaseBody) GetLen() int64 {
	return bb.Len
}

func (bb *BaseBody) GetCfg() map[byte]byte {
	return bb.Cfg
}

func (bb *BaseBody) GetData() map[byte]any {
	return bb.Data
}
