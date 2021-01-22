package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	flag  int
	level Level
	print bool
	mu    *sync.Mutex
	ws    []io.Writer
	skip  int
	ff    LfileFunc
}

func New(flag int, level Level, print bool, ff LfileFunc, ws ...io.Writer) *Logger {
	return &Logger{
		flag:  flag,
		level: level,
		print: print,
		ws:    ws,
		mu:    &sync.Mutex{},
		skip:  4,
		ff:    ff,
	}
}

func (rec *Logger) SetFlag(flag int) *Logger {
	rec.flag = flag
	return rec
}

func (rec *Logger) SetLevel(level Level) *Logger {
	rec.level = level
	return rec
}

func (rec *Logger) SetPrint(print bool) *Logger {
	rec.print = print
	return rec
}

func (rec *Logger) SetLfileFunc(ff LfileFunc) *Logger {
	rec.ff = ff
	return rec
}

func (rec *Logger) SetWriters(ws ...io.Writer) *Logger {
	rec.ws = ws
	return rec
}

func (rec *Logger) UpdateConfig(flag int, level Level, print bool, trimPrefix string, ws ...io.Writer) {
	rec.SetFlag(flag)
	rec.SetLevel(level)
	rec.SetPrint(print)
	rec.SetWriters(ws...)
}

const (
	Ltime = 1 << iota
	Lfile
	Lline
	LstdFlag = Ltime | Lfile | Lline
)

func (rec *Logger) genLog(level string, msg ...interface{}) log {
	var log log

	if rec.flag&Ltime != 0 {
		log.Ltime = genLtime()
	}
	if rec.flag&Lfile != 0 {
		log.Lfile = rec.ff(rec.skip)
	}
	if rec.flag&Lline != 0 {
		log.Lline = genLline(rec.skip)
	}
	log.Level = level
	log.Msg = genMsg(msg...)

	return log
}

func write(w io.Writer, p []byte) (int, error) {
	if w != nil {
		return w.Write(p)
	}
	return 0, nil
}

func (rec *Logger) log(level Level, msg ...interface{}) {
	if rec.level.contains(level) {
		rec.mu.Lock()
		defer rec.mu.Unlock()

		log := rec.genLog(level.string(), msg...).string()
		if rec.print {
			c := level.color()
			c.printLog(log)
		}
		for _, w := range rec.ws {
			write(w, []byte(log))
		}
		if level == LPanic {
			panic(msg)
		}
		if level == LFatal {
			os.Exit(1)
		}
	}
}

func (rec *Logger) Trace(msg ...interface{}) {
	rec.log(LTrace, msg...)
}

func (rec *Logger) Debug(msg ...interface{}) {
	rec.log(LDebug, msg...)
}

func (rec *Logger) Info(msg ...interface{}) {
	rec.log(LInfo, msg...)
}

func (rec *Logger) Warn(msg ...interface{}) {
	rec.log(LWarn, msg...)
}

func (rec *Logger) Error(msg ...interface{}) {
	rec.log(LError, msg...)
}

func (rec *Logger) Panic(msg ...interface{}) {
	rec.log(LPanic, msg...)
}

func (rec *Logger) Fatal(msg ...interface{}) {
	rec.log(LFatal, msg...)
}

func (rec *Logger) HTTP(msg ...interface{}) {
	rec.log(LHTTP, msg...)
}

func (rec *Logger) logf(level Level, format string, msg ...interface{}) {
	if rec.level.contains(level) {
		rec.mu.Lock()
		defer rec.mu.Unlock()

		s := fmt.Sprintf(format, msg...)
		log := rec.genLog(level.string(), s).string()
		if rec.print {
			c := level.color()
			c.printLog(log)
		}
		for _, w := range rec.ws {
			write(w, []byte(log))
		}
		if level == LPanic {
			panic(msg)
		}
		if level == LFatal {
			os.Exit(1)
		}
	}
}

func (rec *Logger) Tracef(format string, msg ...interface{}) {
	rec.logf(LTrace, format, msg...)
}

func (rec *Logger) Debugf(format string, msg ...interface{}) {
	rec.logf(LDebug, format, msg...)
}

func (rec *Logger) Infof(format string, msg ...interface{}) {
	rec.logf(LInfo, format, msg...)
}

func (rec *Logger) Warnf(format string, msg ...interface{}) {
	rec.logf(LWarn, format, msg...)
}

func (rec *Logger) Errorf(format string, msg ...interface{}) {
	rec.logf(LError, format, msg...)
}

func (rec *Logger) Panicf(format string, msg ...interface{}) {
	rec.logf(LPanic, format, msg...)
}

func (rec *Logger) Fatalf(format string, msg ...interface{}) {
	rec.logf(LFatal, format, msg...)
}

func (rec *Logger) HTTPf(format string, msg ...interface{}) {
	rec.logf(LHTTP, format, msg...)
}
