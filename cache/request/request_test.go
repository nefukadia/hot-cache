package request

import (
	"testing"
)

func TestReq(t *testing.T) {
	c := NewReq()
	c.Len = 2
	c.DataBytes = make([]byte, 2)
	c.DataBytes[0] = 0x81
	c.DataBytes[1] = 0x01

	if c.Valid() {
		t.Error("error")
		return
	}
	c.CfgInteger = 0x80
	if c.Valid() {
		t.Error("error")
		return
	}
	c.CfgInteger = ((c.CfgInteger << 8) + 0x83) | (1 << 56)
	if !c.Valid() {
		t.Error("error")
		return
	}

	bytes, err := c.ToBytes()
	if err != nil {
		t.Error(err)
		return
	}

	tmp := c
	c.Len, c.CfgInteger, c.DataBytes = 0, 0, nil
	err = c.Setup(bytes)
	if err != nil {
		t.Error(err)
		return
	}

	if c.Len != tmp.Len || c.CfgInteger != tmp.CfgInteger {
		t.Error("error")
		return
	}
}
