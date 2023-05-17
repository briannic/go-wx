package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

type weatherDataFormatter interface {
	format() string
}

type weatherDataParser interface {
	parseField([]byte) int
}

type weatherDataParserFormatter interface {
	weatherDataFormatter
	weatherDataParser
}

type ApiResults struct {
	data     []weatherDataFormatter
	response []byte
	checksum byte
	length   int
}

func (a *ApiResults) Parse() {
	for i := 5; i < a.length+1; i++ {
		field, found := ApiDefs[a.response[i]]
		if !found {
			fmt.Printf("Did not find field '% X', parsing stopped\n", a.response[i])
			break
		}

		offset := field.parseField(a.response[i:])
		a.data = append(a.data, field)
		i += offset
	}
	return
}

func (r *ApiResults) Display() {
	for i := 0; i < len(r.data); i++ {
		fmt.Printf("%v\n", r.data[i].format())
	}
	fmt.Printf("%v\n", time.Now().Format(time.RFC850))
	fmt.Printf("------\n\n")
}

type field struct {
	id     byte
	label  string
	value  int
	offset int
}

func (d *field) parseField(data []byte) int {
	value := 0
	x := data[1 : d.offset+1]

	switch {
	case d.offset == 1:
		value = int(byte(x[0]))
	case d.offset == 2:
		value = int(binary.BigEndian.Uint16(x))
	case d.offset == 4:
		value = int(binary.BigEndian.Uint32(x))
	case d.offset == 8:
		value = int(binary.BigEndian.Uint64(x))
	}
	d.value = value

	return d.offset
}

type temperatureData struct {
	field
}

func (t temperatureData) format() string {
	unit := "\u00B0"
	value := convertCtoF(float64(t.value / 10))
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type percentageData struct {
	field
}

func (t percentageData) format() string {
	unit := "%"
	return fmt.Sprintf("[%v]%v\t%v%v", t.id, t.label, t.value, unit)
}

type pressureData struct {
	field
}

func (t pressureData) format() string {
	value := convertHpaToInhg(float64(t.value / 10))
	unit := " inhg"
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type directionData struct {
	field
}

func (t directionData) format() string {
	value := float64(t.value)
	unit, _ := lookupCardinalDirection(value)
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type velocityData struct {
	field
}

func (t velocityData) format() string {
	value := convertMsToMph(float64(t.value) / 10)
	unit := " mph"
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type lightData struct {
	field
}

func (t lightData) format() string {
	value := convertLuxToWm2(float64(t.value) / 10)
	unit := " w/m^2"
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type uvData struct {
	field
}

func (t uvData) format() string {
	value := float64(t.value)
	unit := " micro w/m^2"
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type uviData struct {
	field
}

func (t uviData) format() string {
	value := float64(t.value)
	unit := ""
	return fmt.Sprintf("[%v]%v\t%.1f%v", t.id, t.label, value, unit)
}

type rainAmountData struct {
	field
}

func (t rainAmountData) format() string {
	value := convertMmToIn(float64(t.value) / 10)
	unit := " in"
	return fmt.Sprintf("[%v]%v\t%.2f%v", t.id, t.label, value, unit)
}

type rainRateData struct {
	field
}

func (t rainRateData) format() string {
	value := convertMmToIn(float64(t.value) / 10)
	unit := " in/hr"
	return fmt.Sprintf("[%v]%v\t%.2f%v", t.id, t.label, value, unit)
}

var ApiDefs = map[byte]weatherDataParserFormatter{
	0x01: &temperatureData{field{id: 0x01, label: "INTEMP", offset: 2}},
	0x02: &temperatureData{field{id: 0x02, label: "OUTTEMP", offset: 2}},
	0x03: &temperatureData{field{id: 0x03, label: "DEWPOINT", offset: 2}},
	0x04: &temperatureData{field{id: 0x04, label: "WINDCHILL", offset: 2}},
	0x05: &temperatureData{field{id: 0x05, label: "HEATINDEX", offset: 2}},
	0x06: &percentageData{field{id: 0x06, label: "INHUMI", offset: 1}},
	0x07: &percentageData{field{id: 0x07, label: "OUTHUMI", offset: 1}},
	0x08: &pressureData{field{id: 0x08, label: "ABSBARO", offset: 2}},
	0x09: &pressureData{field{id: 0x09, label: "RELBARO", offset: 2}},
	0x0A: &directionData{field{id: 0x0A, label: "WINDDIRECTION", offset: 2}},
	0x0B: &velocityData{field{id: 0x0B, label: "WINDSPEED", offset: 2}},
	0x0C: &velocityData{field{id: 0x0C, label: "GUSTSPEED", offset: 2}},
	0x15: &lightData{field{id: 0x15, label: "LIGHT", offset: 4}},
	0x16: &uvData{field{id: 0x16, label: "UV", offset: 2}},
	0x17: &uviData{field{id: 0x17, label: "UVI", offset: 1}},
	0x19: &velocityData{field{id: 0x19, label: "DAILYWINDMAX", offset: 2}},
	0x80: &rainRateData{field{id: 0x81, label: "PIEZO_RAIN_RATE", offset: 2}},
	0x81: &rainAmountData{field{id: 0x81, label: "PIEZO_EVENT_RAIN", offset: 2}},
	0x82: &rainAmountData{field{id: 0x82, label: "PIEZO_HOURLY_RAIN", offset: 2}},
	0x83: &rainAmountData{field{id: 0x83, label: "PIEZO_DAILY_RAIN", offset: 4}},
	0x84: &rainAmountData{field{id: 0x84, label: "PIEZO_WEEKLY_RAIN", offset: 4}},
	0x85: &rainAmountData{field{id: 0x85, label: "PIEZO_MONTHLY_RAIN", offset: 4}},
	0x86: &rainAmountData{field{id: 0x86, label: "PIEZO_YEARLY_RAIN", offset: 4}},
	0x87: &rainAmountData{field{id: 0x87, label: "PIEZO_GAIN_10", offset: 2 * 10}},
	0x88: &rainAmountData{field{id: 0x88, label: "PIEZO_RST_RAINTIME", offset: 3}},
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

func convertMmToIn(h float64) float64 {
	return h / 25.4
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

func parseResponse(rsp []byte) (ApiResults, error) {
	rspLength := int(binary.BigEndian.Uint16(rsp[3:5]))
	rspChecksum := rsp[rspLength+1]
	results := ApiResults{response: rsp, length: rspLength, checksum: rspChecksum}

	body := rsp[2 : results.length+1]
	if calcChecksum(body) != rspChecksum {
		return results, errors.New("Checksum mismatch")
	}

	results.Parse()

	return results, nil
}
