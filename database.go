package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type WeatherDB struct {
	db *sql.DB
}

func (w *WeatherDB) Insert(a *ApiResults) {
	//db, err := sql.Open("sqlite3", "data/go_wx.db")
	//checkErr(err)

	//_, err := db.Prepare(
	//	"INSERT INTO wx_data(" +
	//		"id, time, INTEMP, OUTTEMP," +
	//		"DEWPOINT, WINDCHILL, HEATINDEX," +
	//		"INHUMI, OUTHUMI, ABSBARO," +
	//		"RELBARO, WINDDIRECTION, WINDSPEED," +
	//		"GUSTSPEED, LIGHT, UV, UVI," +
	//		"DAILYWINDMAX, PIEZO_RAIN_RATE," +
	//		"PIEZO_EVENT_RAIN, PIEZO_HOURLY_RAIN," +
	//		"PIEZO_DAILY_RAIN, PIEZO_WEEKLY_RAIN," +
	//		"PIEZO_MONTHLY_RAIN, PIEZO_YEARLY_RAIN," +
	//		"PIEZO_GAIN_10, PIEZO_RST_RAINTIME)" +
	//		"values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	//checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
