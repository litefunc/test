package main

import (
	"cloud/lib/logger"
	"errors"
	"fmt"
)

func main() {
	e := errors.New("err")
	var err error
	err = fmt.Errorf(".%v. %v .%v.", e, e, e)
	logger.Debug(err)
	logger.Debug(errors.Unwrap(err))
}
