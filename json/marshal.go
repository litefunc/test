package main

import (
	"VodoPlay/logger"
	"encoding/json"
	"math"
)

func main() {
	var n interface{}
	var err error
	var by []byte

	by, err = json.Marshal(n)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	by, err = json.Marshal(1)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	var ch chan int
	by, err = json.Marshal(ch)
	logger.Debugf(`%s`, by)
	logger.Error(err)

	by, err = json.Marshal(math.Inf(1))
	logger.Debugf(`%s`, by)
	logger.Error(err)

	var p *struct{}
	by, err = json.Marshal(p)
	logger.Debugf(`%s`, by)
	logger.Error(err)
}
