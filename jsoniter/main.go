package main

import (
	"cloud/lib/logger"

	jsoniter "github.com/json-iterator/go"
)

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {

	// group := ColorGroup{
	// 	ID:     1,
	// 	Name:   "Reds",
	// 	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	// }
	b, err := jsoniter.Marshal(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(b))
	logger.Debug(jsoniter.Get(b, "Colors", 0).ToString())

	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	logger.Debug(string(val))
	logger.Debug(jsoniter.Get(val, "Colors", 0).ToString())
}
