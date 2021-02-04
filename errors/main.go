package main

import (
	"errors"
	"fmt"
	"test/errors/internal"
	"test/logger"
)

func main() {
	e := errors.New("err")
	var err error
	err = fmt.Errorf(".%v. %v .%v.", e, e, e)
	logger.Debug(err)
	logger.Debug(errors.Unwrap(err))
	err = fmt.Errorf(".%w. %w .%w.", e, e, e)
	logger.Debug(err)
	logger.Debug(errors.Unwrap(err))
	err = fmt.Errorf(".%w", e)
	logger.Debug(err)
	logger.Debug(errors.Unwrap(err))

	internal.Unwrap()
	internal.Is()
	internal.As()
}
