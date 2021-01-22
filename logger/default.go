package logger

import (
	"os"
	"sync"
)

var std = &Logger{
	flag:  LstdFlag,
	level: LstdLevel,
	print: true,
	mu:    &sync.Mutex{},
	skip:  5,
	ff:    RelLfile,
}

func Default() *Logger {
	return std
}

func Trace(msg ...interface{}) {
	std.Trace(msg...)
}

func Debug(msg ...interface{}) {
	std.Debug(msg...)
}

func Info(msg ...interface{}) {
	std.Info(msg...)
}

func Warn(msg ...interface{}) {
	std.Warn(msg...)
}

func Error(msg ...interface{}) {
	std.Error(msg...)
}

func Panic(msg ...interface{}) {
	std.Panic(msg...)
}

func Fatal(msg ...interface{}) {
	std.Fatal(msg...)
}

func HTTP(msg ...interface{}) {
	std.HTTP(msg...)
}

func Tracef(format string, msg ...interface{}) {
	std.Tracef(format, msg...)
}

func Debugf(format string, msg ...interface{}) {
	std.Debugf(format, msg...)
}

func Infof(format string, msg ...interface{}) {
	std.Infof(format, msg...)
}

func Warnf(format string, msg ...interface{}) {
	std.Warnf(format, msg...)
}

func Errorf(format string, msg ...interface{}) {
	std.Errorf(format, msg...)
}

func Panicf(format string, msg ...interface{}) {
	std.Panicf(format, msg...)
}

func Fatalf(format string, msg ...interface{}) {
	std.Fatalf(format, msg...)
}

func HTTPf(format string, msg ...interface{}) {
	std.HTTPf(format, msg...)
}

func Exit(err error) {
	if err != nil {
		msgs := []interface{}{err.Error()}
		logs := std.genLog("Exit", msgs...)
		red := newColor(FgBrRed)
		red.printLog(logs.string())
		os.Exit(1)
	}
}

func LogErr(err error) {
	if err != nil {
		msgs := []interface{}{err.Error()}
		logs := std.genLog("Error", msgs...)
		red := newColor(FgBrRed)
		red.printLog(logs.string())
	}
}
