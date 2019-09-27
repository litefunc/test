package pgsql

import (
	"log"
	"testing"
)

type Tester struct {
	t *testing.T
}

func NewTester(t *testing.T) *Tester {
	return &Tester{t}
}

type Runner interface {
	Run() error
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func Run(i Runner) error {
	return i.Run()
}

func Scan(i Scanner, dest ...interface{}) error {
	return i.Scan(dest...)
}

func (t Tester) Run(i Runner) {
	if err := i.Run(); err != nil {
		t.t.Error(err)
		log.Panic(err)
	}
}

func (t Tester) Scan(i Scanner, dest ...interface{}) {
	if err := i.Scan(dest...); err != nil {
		t.t.Error(err)
		log.Panic(err)
	}
}
