package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type log struct {
	Ltime string
	Lfile string
	Lline int
	Level string
	Msg   string
}

func genLtime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

type LfileFunc func(skip int) string

func TrimPrefixLfileFunc(trimPrefix string) LfileFunc {
	f := func(skip int) string {
		_, file, _, _ := runtime.Caller(skip)
		return strings.TrimPrefix(file, trimPrefix)
	}
	return f
}

func AbsLfile(skip int) string {
	_, file, _, _ := runtime.Caller(skip)
	return file
}

func RelLfile(skip int) string {
	_, file, _, _ := runtime.Caller(skip)
	pwd, _ := os.Getwd()
	file = strings.TrimPrefix(file, pwd+"/")
	return file
}

func genLline(skip int) int {
	_, _, line, _ := runtime.Caller(skip)
	return line
}

func genMsg(msg ...interface{}) string {
	var msgs []string
	for _, v := range msg {
		msgs = append(msgs, fmt.Sprintf("%+v", v))
	}
	return strings.Join(msgs, " ")
}

func (log log) string() string {
	var msgs []string

	if log.Ltime != "" {
		msgs = append(msgs, log.Ltime)
	}

	if log.Level != "" {
		n := 5 - len(log.Level)
		if n < 0 {
			n = 0
		}
		msgs = append(msgs, fmt.Sprintf("[%s]%s", log.Level, strings.Repeat(" ", n)))
	}

	s := ""
	if log.Lfile != "" {
		s += fmt.Sprintf("%s:", log.Lfile)
	}
	if log.Lline != 0 {
		s += fmt.Sprintf("%d:", log.Lline)
	}
	if s != "" {
		msgs = append(msgs, s)
	}

	if log.Msg != "" {
		msgs = append(msgs, log.Msg)
	}

	return strings.Join(msgs, " ") + "\n"
}
