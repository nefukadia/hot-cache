package response

import (
	"testing"
)

func TestResp(t *testing.T) {
	c := NewResp()
	c.Len = 9
	c.DataBytes = make([]byte, 9)
	c.DataBytes[0] = 0x00
	c.DataBytes[1] = 0x00
	c.DataBytes[2] = 0x00
	c.DataBytes[3] = 0x00
	c.DataBytes[4] = 0x00
	c.DataBytes[5] = 0x00
	c.DataBytes[6] = 0x00
	c.DataBytes[7] = 0x01
	c.DataBytes[8] = "a"[0]
	if c.Valid() {
		t.Error("error")
		return
	}
	c.CfgInteger = 0x60
	if c.Valid() {
		t.Error("error")
		return
	}
	c.CfgInteger = ((c.CfgInteger << 8) + 0x96) | (uint64(255) << 56)
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

	err = c.SetupValue(0, 1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	testValue := "s1"
	testInfo := "s2"
	err = c.SetupValue(0, 2, &testValue, &testInfo)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetupError(3, 0, "error")
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetupNotFound(4, 5)
	if err != nil {
		t.Fatal(err)
	}
}
