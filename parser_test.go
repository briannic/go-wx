package main

import (
	"reflect"
	"testing"
)

func TestParseResponse(t *testing.T) {
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

	want := ApiResults{
		response: input,
		data: []weatherDataFormatter{
			&temperatureData{field{id: 0x01, label: "INTEMP", value: 223}},
			&percentageData{field{id: 0x06, label: "INHUMI", value: 35}},
			&pressureData{field{id: 0x08, label: "ABSBARO", value: 9768}},
			&pressureData{field{id: 0x09, label: "RELBARO", value: 9768}},
			&temperatureData{field{id: 0x02, label: "OUTTEMP", value: 210}},
			&percentageData{field{id: 0x07, label: "OUTHUMI", value: 42}},
			&directionData{field{id: 0x0A, label: "WINDDIRECTION", value: 65}},
			&velocityData{field{id: 0x0B, label: "WINDSPEED", value: 0}},
			&velocityData{field{id: 0x0C, label: "GUSTSPEED", value: 6}},
			&lightData{field{id: 0x15, label: "LIGHT", value: 0}},
			&uvData{field{id: 0x16, label: "UV", value: 0}},
			&uviData{field{id: 0x17, label: "UVI", value: 0}},
			&velocityData{field{id: 0x19, label: "DAILYWINDMAX", value: 16}},
		},
		checksum: 223,
		length:   42,
	}

	got, error := parseResponse(input)

	if error != nil {
		t.Errorf("Error: got %v, wanted %v", error, nil)
	}

	if !reflect.DeepEqual(got.response, want.response) {
		t.Errorf("Response: got %v, wanted %v", got.response, want.response)
	}

	if got.checksum != want.checksum {
		t.Errorf("Checksum: got %v, wanted %v", got.checksum, want.checksum)
	}

	if got.length != want.length {
		t.Errorf("Length: got %v, wanted %v", got.length, want.length)
	}

	for i, v := range got.data {
		got_format := v.format()
		want_format := want.data[i].format()

		if !reflect.DeepEqual(got_format, want_format) {
			t.Errorf("Format: got %v, wanted %v", got_format, want_format)
		}

	}
}

func TestCalcChecksum(t *testing.T) {
	want := byte(42)
	input := []byte("\x27\x03")
	got := calcChecksum(input)

	if want != got {
		t.Errorf("got %x, wanted %x", got, want)
	}
}
