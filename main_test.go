package main

import (
	"reflect"
	"testing"
)

func TestCalcChecksum(t *testing.T) {
	want := byte(42)
	input := []byte("\x27\x03")
	got := calcChecksum(input)

	if want != got {
		t.Errorf("got %x, wanted %x", got, want)
	}
}

func TestCreateMsg(t *testing.T) {
	want := []byte{0xFF, 0xFF, 0x80, 0x04, 0x01, 0x85}
	got := createMsg([]byte{0x80}, []byte{0x01})

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %x, wanted %x", got, want)
	}
}
