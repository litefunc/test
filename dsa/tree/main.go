package main

import (
	"encoding/json"
	"test/dsa/tree/internal"
	"test/logger"
)

func log(v interface{}) string {
	by, _ := json.MarshalIndent(v, "", "\t")
	return string(by)
}

func main() {

	n := internal.NewNode(true)
	n1 := internal.NewNode(true)
	n2 := internal.NewNode(false)
	n3 := internal.NewNode(false)

	n.Set("1", n1)
	n.Set("3", n3)
	n1.Set("2", n2)
	logger.Debug(log(n))

	n5 := n.Get("1", "2")
	logger.Debug(log(n5))
	n5.B = true

	logger.Debug(log(n))

	n5.Get("1", "2")
	n.Get("4")

}
