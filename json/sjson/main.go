package main

import (
	"cloud/lib/logger"

	"github.com/tidwall/sjson"
)

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main() {
	value, err := sjson.Set(json, "name.last", "Anderson")
	logger.Debug(value, err)

	value, err = sjson.Set(json, "name.last1", []int{1, 2, 3})
	logger.Debug(value, err)
}
