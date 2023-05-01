package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type ApiDef struct {
	label  string
	offset int
}

type ApiResults struct {
	data []Datapoint
}

type Datapoint struct {
	id    byte
	label string
	value int
}

var ApiDefs = map[byte]ApiDef{
	0x01: {"INTEMP", 2},
	0x02: {"OUTTEMP", 2},
	0x03: {"DEWPOINT", 2},
	0x04: {"WINDCHILL", 2},
	0x05: {"HEATINDEX", 2},
	0x06: {"INHUMI", 1},
	0x07: {"OUTHUMI", 1},
	0x08: {"ABSBARO", 2},
	0x09: {"RELBARO", 2},
	0x0A: {"WINDDIRECTION", 2},
	0x0B: {"WINDSPEED", 2},
	0x0C: {"GUSTSPEED", 2},
	0x0D: {"RAINEVENT", 2},
	0x0E: {"RAINRATE", 2},
	0x0F: {"RAINHOUR", 2},
	0x10: {"RAINDAY", 2},
	0x11: {"RAINWEEK", 2},
	0x12: {"RAINMONTH", 4},
	0x13: {"RAINYEAR", 4},
	0x14: {"RAINTOTALS", 4},
	0x15: {"LIGHT", 4},
	0x16: {"UV", 2},
	0x17: {"UVI", 1},
	0x18: {"TIME", 6},
	0x19: {"DAILYWINDMAX", 2},
	0x1A: {"TEMP1", 2},
	0x1B: {"TEMP2", 2},
	0x1C: {"TEMP3", 2},
	0x1D: {"TEMP4", 2},
	0x1E: {"TEMP5", 2},
	0x1F: {"TEMP6", 2},
	0x20: {"TEMP7", 2},
	0x21: {"TEMP8", 2},
	0x22: {"HUM1", 2},
	0x23: {"HUM2", 2},
	0x24: {"HUM3", 2},
	0x25: {"HUM4", 2},
	0x26: {"HUM5", 2},
	0x27: {"HUM6", 2},
	0x28: {"HUM7", 2},
	0x29: {"HUM8", 2},
	0x4C: {"LOWBAT", 16},
	0x80: {"PIEZO_RAIN_RATE", 2},
	0x81: {"PIEZO_EVENT_RAIN", 2},
	0x82: {"PIEZO_HOURLY_RAIN", 2},
	0x83: {"PIEZO_DAILY_RAIN", 4},
	0x84: {"PIEZO_WEEKLY_RAIN", 4},
	0x85: {"PIEZO_MONTHLY_RAIN", 4},
	0x86: {"PIEZO_YEARLY_RAIN", 4},
	0x87: {"PIEZO_GAIN_10", 2 * 10},
	0x88: {"PIEZO_RST_RAINTIME", 3},
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var tempFields = []string{
	"INTEMP",
	"OUTTEMP",
}

func convertCtoF(c float64) float64 {
	f := c / 10
	return (f * 1.8) + 32
}

func (d *Datapoint) Transform() {
	value := float64(d.value)

	switch {
	case stringInSlice(d.label, tempFields):
		value = convertCtoF(value)
	}

	fmt.Printf("[%v]%v\t%v\n", d.id, d.label, value)
}

func (r *ApiResults) Display() {
	for i := 0; i < len(r.data); i++ {
		r.data[i].Transform()
	}
}

func parseResponse(rsp []byte) (ApiResults, error) {
	results := ApiResults{}
	rspLength := int(binary.BigEndian.Uint16(rsp[3:5]))
	rspChecksum := rsp[rspLength+1]
	body := rsp[2 : rspLength+1]

	if calcChecksum(body) != rspChecksum {
		return results, errors.New("Checksum mismatch")
	}

	for i := 5; i < rspLength+1; i++ {
		def := ApiDefs[rsp[i]]

		val := 0
		x := rsp[i+1 : i+def.offset+1]

		switch {
		case def.offset == 1:
			val = int(byte(x[0]))
		case def.offset == 2:
			val = int(binary.BigEndian.Uint16(x))
		case def.offset == 4:
			val = int(binary.BigEndian.Uint32(x))
		case def.offset == 8:
			val = int(binary.BigEndian.Uint64(x))
		}

		results.data = append(results.data, Datapoint{rsp[i], def.label, val})
		i += def.offset
	}
	return results, nil
}
