package main

import (
	"cloud/lib/logger"

	"go.mongodb.org/mongo-driver/bson"
)

type D []E

type E struct {
	Key   string
	Value interface{}
}

func main() {
	d := D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
	m := bson.M{"name": "pi", "value": 3.14159}
	logger.Debug(d)
	logger.Debug(m["name"])

	// var m1 map[string]interface{}
	m1 := map[string]interface{}{"name": "pi", "value": 3.14159}
	m1["name"] = "name1"
	logger.Debug(m1["name"])
	f(m)
	f1(m1)
}

func f(map[string]interface{}) {}
func f1(bson.M)                {}
