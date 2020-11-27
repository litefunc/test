package main

import (
	"cloud/lib/logger"
	"fmt"

	"github.com/tidwall/sjson"
)

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main() {

	value, err := sjson.Set("", "name.last", "Anderson")
	logger.Debug(value, err)

	value, err = sjson.Set(json, "name.last", "Anderson")
	logger.Debug(value, err)

	value, err = sjson.Set(json, "name.last1", []int{1, 2, 3})
	logger.Debug(value, err)
	value, err = sjson.Set(json, "name1.last1", []int{1, 2, 3})
	logger.Debug(value, err)

	value, err = sjson.Delete(json, "name.last")
	logger.Debug(value, err)
	value, err = sjson.Delete(json, "name")
	logger.Debug(value, err)

	value, err = sjson.Set("", fmt.Sprintf("auths.%s.auth", "noovo-dock.ddns.net:5000"), "123")
	logger.Debug(value, err)

	value, err = sjson.Set("", "auths", "noovo-dock.ddns.net:5000")
	logger.Debug(value, err)
	value, err = sjson.Set("", "auths.noovo-dock.ddns.net:5000", "a")
	logger.Debug(value, err)
}
