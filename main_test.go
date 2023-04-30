package main

import (
	"testing"
)

func TestCalcChecksum(t *testing.T) {

	want := byte(42)
	input := [][]byte{[]byte("\x27\x03")}
	got := calcChecksum(input)

	if want != got {
		t.Errorf("got %x, wanted %x", got, want)
	}
}
