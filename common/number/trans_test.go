package number

import (
	"testing"
)

func TestLowByteToOther(t *testing.T) {
	x := 255
	target := LowByteToOther[int, int](x)
	if target != x {
		t.Fatal("error")
	}

	x = 257
	target = LowByteToOther[int, int](x)
	if target != x-256 {
		t.Fatal("error")
	}
}

func TestToStringPtr(t *testing.T) {
	if *ToStringPtr(101) != "101" {
		t.Fatal("error")
	}
}
