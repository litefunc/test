package internal

import (
	"errors"
	"fmt"
	"reflect"
	"test/logger"
)

var (
	err1 = Err1{
		Err: err2,
	}
	err2 = Err2{
		Err: err3,
	}
	err3 = Err3{
		Err: errors.New("Err 3"),
	}
	err4 = errors.New("err4")
)

type Err1 struct {
	Err error
}

func (rec Err1) Error() string {
	return fmt.Sprintf(`%d:%v`, 1, rec.Err)
}

func (rec Err1) Unwrap() error {
	return rec.Err
}

func (rec Err1) Is(err error) bool {
	return reflect.DeepEqual(rec, err)
}

type Err2 struct {
	Err error
}

func (rec Err2) Error() string {
	// return fmt.Sprintf(`%d:%w`, 2, rec.Err)
	// return fmt.Sprintf(`%d:%s`, 2, rec.Err)
	return fmt.Sprintf(`%d:%v`, 2, rec.Err)
}

func (rec Err2) Unwrap() error {
	return rec.Err
}

type Err3 struct {
	Err error
}

func (rec Err3) Error() string {
	return fmt.Sprintf(`%d:%v`, 3, rec.Err)
}

func (rec Err3) Unwrap() error {
	return rec.Err
}

func (rec Err3) As(i interface{}) bool {
	switch i.(type) {

	case *Err4, *Err6:
		return true

	default:
		return false
	}
}

type Err4 struct {
	Err error
}

func (rec Err4) Error() string {
	return fmt.Sprintf(`%d:%v`, 3, rec.Err)
}

type Err5 struct{}

func (rec Err5) Error() string {
	return ""
}

type Err6 struct{}

func genErr() error {
	// return err1
	return fmt.Errorf(`0 %w`, err1)
}

func genErr1() error {

	err3 := Err3{
		Err: errors.New("Err 3"),
	}

	err2 := Err2{
		Err: err3,
	}

	err1 := Err1{
		Err: err2,
	}

	return fmt.Errorf(`%w`, err1)
}

func Unwrap() {
	err := genErr()
	unwrap(err)
}

func unwrap(err error) {
	logger.Info(err)
	n := 1
	for {
		err = errors.Unwrap(err)
		logger.Info(n, err)
		if err == nil {
			break
		}
		n++
	}
}

func Is() {
	err := genErr()
	logger.Debug(errors.Is(err, err1))
	logger.Debug(errors.Is(err, err2))
	logger.Debug(errors.Is(err, err3))
	logger.Debug(errors.Is(err, err4))

	err = genErr1()
	logger.Debug(errors.Is(err, err1))
	logger.Debug(errors.Is(err, err2))
	logger.Debug(errors.Is(err, err3))
	logger.Debug(errors.Is(err, err4))
}

func As() {
	err := genErr()
	logger.Debug(errors.As(err, &Err1{}))
	logger.Debug(errors.As(err, &Err2{}))
	logger.Debug(errors.As(err, &Err3{}))
	logger.Debug(errors.As(err, &Err4{}))
	logger.Debug(errors.As(err, &Err5{}))
	logger.Debug(errors.As(err, &Err6{}))
}
