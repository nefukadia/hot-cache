package cache

import (
	"unsafe"
)

const SizeByte = int(unsafe.Sizeof(0))
const OptionByte = int(unsafe.Sizeof(Option(0)))

type Option byte
type RepCode byte

type Res struct {
	Op   Option
	Data []byte
}

type Rep struct {
	Code RepCode
	Data []byte
}

func (rep *Rep) Len() int {
	return int(unsafe.Sizeof(rep.Code)) + len(rep.Data)
}

type Cache struct {
}

func New() *Cache {
	return &Cache{}
}

func (cache *Cache) Solve(res Res) Rep {
	// todo
	return Rep{}
}
