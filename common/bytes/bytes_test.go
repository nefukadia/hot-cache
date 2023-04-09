package bytes

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	var data []byte
	for _, b := range data {
		fmt.Printf("%08b ", b)
	}
	data = VLEMarshal(data)
	for _, b := range data {
		fmt.Printf("%08b ", b)
	}
}

func TestUnmarshal(t *testing.T) {
	var data []byte
	data, err := VLEUnmarshal(data)
	if err != nil {
		panic(err)
	}
}
