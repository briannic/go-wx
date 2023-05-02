package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

type ApiDef struct {
	label  string
	offset int
}

type ApiResults struct {
	data      []Datapoint
	createdAt time.Time
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

var temperatureFields = []string{
	"INTEMP",
	"OUTTEMP",
}

var pressureFields = []string{
	"ABSBARO",
	"RELBARO",
}

var percentageFields = []string{
	"INHUMI",
	"OUTHUMI",
}

var velocityFields = []string{
	"WINDSPEED",
	"GUSTSPEED",
	"DAILYWINDMAX",
}

var lightFields = []string{
	"LIGHT",
}

var directionFields = []string{
	"WINDDIRECTION",
}

func convertCtoF(c float64) float64 {
	return (c * 1.8) + 32
}

func convertHpaToInhg(h float64) float64 {
	return h / 33.8638
}

func convertMsToMph(h float64) float64 {
	return h * 2.237
}

func convertLuxToWm2(h float64) float64 {
	return h * 0.0079
}

func lookupCardinalDirection(degree float64) (string, error) {
	if degree < 0 || degree > 360 {
		return "", errors.New("invalid direction, must be 0<=x<=360")
	}

	dir := ""
	switch {
	case degree >= 348 || degree < 11:
		dir = "N"
	case degree >= 11 && degree < 33:
		dir = "NNE"
	case degree >= 33 && degree < 56:
		dir = "NE"
	case degree >= 56 && degree < 78:
		dir = "ENE"
	case degree >= 78 && degree < 101:
		dir = "E"
	case degree >= 101 && degree < 123:
		dir = "ESE"
	case degree >= 123 && degree < 146:
		dir = "SE"
	case degree >= 146 && degree < 168:
		dir = "SSE"
	case degree >= 168 && degree < 191:
		dir = "S"
	case degree >= 191 && degree < 213:
		dir = "SSW"
	case degree >= 213 && degree < 236:
		dir = "SW"
	case degree >= 236 && degree < 258:
		dir = "WSW"
	case degree >= 258 && degree < 281:
		dir = "W"
	case degree >= 281 && degree < 303:
		dir = "WNW"
	case degree >= 303 && degree < 326:
		dir = "NW"
	case degree >= 326 && degree < 348:
		dir = "NNW"
	}
	return dir, nil

}

func (d *Datapoint) Transform() {
	value := float64(d.value)
	unit := ""

	switch {
	case stringInSlice(d.label, temperatureFields):
		value = convertCtoF(value / 10)
		unit = "\u00B0"
	case stringInSlice(d.label, pressureFields):
		value = convertHpaToInhg(value / 10)
		unit = " inhg"
	case stringInSlice(d.label, percentageFields):
		unit = "%"
	case stringInSlice(d.label, velocityFields):
		value = convertMsToMph(value / 10)
		unit = " mph"
	case stringInSlice(d.label, lightFields):
		value = convertLuxToWm2(value / 10)
		unit = " w/m^2"
	case stringInSlice(d.label, directionFields):
		unit, _ = lookupCardinalDirection(value)
	}

	fmt.Printf("[%v]%v\t%.1f%v\n", d.id, d.label, value, unit)
}

func (r *ApiResults) Display() {
	for i := 0; i < len(r.data); i++ {
		r.data[i].Transform()
	}
	fmt.Printf("%v\n", r.createdAt.Format(time.RFC850))
	fmt.Println("------")
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
		def, found := ApiDefs[rsp[i]]
		if !found {
			fmt.Printf("Did not find field '% X', parsing stopped\n", rsp[i])
			break
		}

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
		results.createdAt = time.Now()
		i += def.offset
	}
	return results, nil
}
