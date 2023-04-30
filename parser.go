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

var ApiDefs = map[string]ApiDef{
	"\x01": {"INTEMP", 2},
	"\x02": {"OUTTEMP", 2},
	"\x03": {"DEWPOINT", 2},
	"\x04": {"WINDCHILL", 2},
	"\x05": {"HEATINDEX", 2},
	"\x06": {"INHUMI", 1},
	"\x07": {"OUTHUMI", 1},
	"\x08": {"ABSBARO", 2},
	"\x09": {"RELBARO", 2},
	"\x0A": {"WINDDIRECTION", 2},
	"\x0B": {"WINDSPEED", 2},
	"\x0C": {"GUSTSPEED", 2},
	"\x0D": {"RAINEVENT", 2},
	"\x0E": {"RAINRATE", 2},
	"\x0F": {"RAINHOUR", 2},
	"\x10": {"RAINDAY", 2},
	"\x11": {"RAINWEEK", 2},
	"\x12": {"RAINMONTH", 4},
	"\x13": {"RAINYEAR", 4},
	"\x14": {"RAINTOTALS", 4},
	"\x15": {"LIGHT", 4},
	"\x16": {"UV", 2},
	"\x17": {"UVI", 1},
	"\x18": {"TIME", 6},
	"\x19": {"DAILYWINDMAX", 2},

	"\x1A": {"TEMP1", 2},
	"\x1B": {"TEMP2", 2},
	"\x1C": {"TEMP3", 2},
	"\x1D": {"TEMP4", 2},
	"\x1E": {"TEMP5", 2},
	"\x1F": {"TEMP6", 2},
	"\x20": {"TEMP7", 2},
	"\x21": {"TEMP8", 2},

	"\x22": {"HUM1", 2},
	"\x23": {"HUM2", 2},
	"\x24": {"HUM3", 2},
	"\x25": {"HUM4", 2},
	"\x26": {"HUM5", 2},
	"\x27": {"HUM6", 2},
	"\x28": {"HUM7", 2},
	"\x29": {"HUM8", 2},

	"\x4C": {"LOWBAT", 16},

	"\x80": {"PIEZO_RAIN_RATE", 2},
	"\x81": {"PIEZO_EVENT_RAIN", 2},
	"\x82": {"PIEZO_HOURLY_RAIN", 2},
	"\x83": {"PIEZO_DAILY_RAIN", 4},
	"\x84": {"PIEZO_WEEKLY_RAIN", 4},
	"\x85": {"PIEZO_MONTHLY_RAIN", 4},
	"\x86": {"PIEZO_YEARLY_RAIN", 4},
	"\x87": {"PIEZO_GAIN_10", 2 * 10},
	"\x88": {"PIEZO_RST_RAINTIME", 3},
}

func transformData(label string, data int) {
}

func parseResponse(rsp []byte) error {
	rspLength := int(binary.BigEndian.Uint16(rsp[3:5]))
	rspChecksum := rsp[rspLength+1]
	body := rsp[2 : rspLength+1]
	if sum(body)%256 != int(rspChecksum) {
		return errors.New("Checksum mismatch")
	}

	for i := 5; i < rspLength+1; i++ {
		def := ApiDefs[string(rsp[i])]

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

		fmt.Printf("%v[%v] %v\n", rsp[i], def.label, val)
		i += def.offset
	}
	return nil
}
