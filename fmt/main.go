package main

import (
	"cloud/lib/logger"
	"errors"
	"fmt"
)

func main() {
	err := fmt.Errorf(`%v, %v`, errors.New("e1"), errors.New("e2"))
	logger.Debug(err)
}
