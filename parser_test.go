package main

import (
	"reflect"
	"testing"
)

func TestParseResponse(t *testing.T) {
	want := ApiResults{
		[]Datapoint{
			{id: 0x01, label: "INTEMP", value: 223},
			{id: 0x06, label: "INHUMI", value: 35},
			{id: 0x08, label: "ABSBARO", value: 9768},
			{id: 0x09, label: "RELBARO", value: 9768},
			{id: 0x02, label: "OUTTEMP", value: 210},
			{id: 0x07, label: "OUTHUMI", value: 42},
			{id: 0x0A, label: "WINDDIRECTION", value: 65},
			{id: 0x0B, label: "WINDSPEED", value: 0},
			{id: 0x0C, label: "GUSTSPEED", value: 6},
			{id: 0x15, label: "LIGHT", value: 0},
			{id: 0x16, label: "UV", value: 0},
			{id: 0x17, label: "UVI", value: 0},
			{id: 0x19, label: "DAILYWINDMAX", value: 16},
		},
	}

	input := []byte("\xff\xff\x27\x00\x2A" +
		"\x01\x00\xDF\x06\x23" +
		"\x08\x26\x28\x09\x26" +
		"\x28\x02\x00\xD2\x07" +
		"\x2A\x0A\x00\x41\x0B" +
		"\x00\x00\x0C\x00\x06" +
		"\x15\x00\x00\x00\x00" +
		"\x16\x00\x00\x17\x00" +
		"\x19\x00\x10\xDF",
	)

	got, error := parseResponse(input)

	if error != nil {
		t.Errorf("got %v, wanted %v", error, nil)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
