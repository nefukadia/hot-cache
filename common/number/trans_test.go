package number

import (
	"testing"
)

func TestLowByteToOther(t *testing.T) {
	x := 255
	target := LowByteToOther[int, int](x)
	if target != x {
		t.Error("error")
	}

	x = 257
	target = LowByteToOther[int, int](x)
	if target != x-256 {
		t.Error("error")
	}
}
