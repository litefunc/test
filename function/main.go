package main

import (
	"test/logger"
	"time"
)

func fun(f func(interface{}) interface{}, o interface{}) interface{} {
	return f(o)
}

func f1(i int) bool {
	return true
}

func f2(i int) {
	i++
}

func main() {
	i := 0
	f2(i)
	logger.Debug(i)
}

type syncType int

const (
	st_thumbnail syncType = 0
	st_time      syncType = 1
	st_unix      syncType = 2
)

type Sync struct {
	data     data
	hasItems bool
	items    map[string]*Sync
}

type data struct {
	id       string
	sync     bool
	st       syncType
	time     time.Time
	unix     int64
	json     string
	filepath string
}

func (rec *Sync) Apply(f func(*data)) {

	if !rec.hasItems {
		f(&rec.data)
		return
	}

	for _, v := range rec.items {
		v.Apply(f)
	}
}

func (rec *Sync) Sync() bool {

	if !rec.hasItems {
		return rec.data.sync
	}

	for _, v := range rec.items {
		if !v.Sync() {
			return false
		}
	}
	return true
}
