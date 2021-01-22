package main

import (
	"fmt"
	"os"
	"sync"
	"test/logger"
	"time"
)

func main() {

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}
	d := logger.Default()
	d.SetWriters(f, f)

	fmt.Println(logger.LHTTP, logger.LstdLevel)
	fmt.Println(logger.LstdFlag)

	logger.Trace("a", 1)
	logger.Debug("a", 1)
	logger.Info("a", 1)
	logger.Warn("a", 1)
	logger.Error("a", 1)
	logger.HTTP("a", 1)

	logger.Tracef(`%s, %d`, "a", 1)
	logger.Debugf(`%s, %d`, "a", 1)
	logger.Infof(`%s, %d`, "a", 1)
	logger.Warnf(`%s, %d`, "a", 1)
	logger.Errorf(`%s, %d`, "a", 1)
	logger.HTTPf(`%s, %d`, "a", 1)

	l := logger.New(logger.LstdFlag, logger.LstdLevel, true, logger.RelLfile, f)

	l.Trace("a", 1)
	l.Debug("a", 1)
	l.Info("a", 1)
	l.Warn("a", 1)
	l.Error("a", 1)
	l.HTTP("a", 1)

	l.Tracef(`%s, %d`, "a", 1)
	l.Debugf(`%s, %d`, "a", 1)
	l.Infof(`%s, %d`, "a", 1)
	l.Warnf(`%s, %d`, "a", 1)
	l.Errorf(`%s, %d`, "a", 1)
	l.HTTPf(`%s, %d`, "a", 1)

	// logger.Panic("a", 1)
	// logger.Panicf(`%s, %d`, "a", 1)
	// logger.Fatal("a", 1)
	// logger.Fatalf(`%s, %d`, "a", 1)

	var wc sync.WaitGroup
	wc.Add(2)
	t := time.NewTicker(time.Second * 5)
	for range t.C {

		go func() {
			logger.Debug(1, time.Now())
		}()
		go func() {
			logger.Debug(2, time.Now())
		}()
		go func() {
			logger.Debug(3, time.Now())
		}()
	}
	wc.Wait()
}
