package logger

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

// Logger represents an active logging object that generates lines ofoutput
type Logger struct {
	flag  int
	level int
	save  bool
	file  string
}

const (
	Ltime = 1 << iota
	Lfile
	Lline
	Ltype
	Ldebug
)

const (
	LTrace = 1 << iota
	LDebug
	LInfo
	LWarn
	LError
	LFatal
	LPanic
	LHTTP
)

func SetLogger(flag, level int, save bool, file string, l *Logger) {
	SetFlags(flag, l)
	SetLevel(level, l)
	SetSave(save, l)
	SetFile(file, l)
}

// SetFlags represent developer can customize logger
func SetFlags(flag int, l *Logger) {
	l.flag = flag
}

func SetLevel(level int, l *Logger) {
	l.level = level
}

func SetSave(save bool, l *Logger) {
	l.save = save
}

func SetFile(file string, l *Logger) {
	l.file = file
}

func printLog(c *color, log string) {
	c.println(fmt.Sprintf("%v", log))
}

func writeToFile(filename string, log string) {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	if _, err := f.WriteString(log); err != nil {
		fmt.Println(err)
	}
}

func SetStderr(f string) {

	logFile, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	redirectStderr(logFile)
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Printf("Failed to redirect stderr to file: %v", err)
	}
}

func (rec Logger) Trace(msg ...interface{}) {
	logs := genLog("Trace", msg...)

	if rec.level&LTrace != 0 {
		log := GenLog(logs)
		color := newColor(FgBrBlue)
		printLog(color, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Debug(msg ...interface{}) {
	logs := genLog("Debug", msg...)

	if rec.level&LDebug != 0 {
		log := GenLog(logs)
		color := newColor(FgBrCyan)
		printLog(color, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Info(msg ...interface{}) {
	logs := genLog("Info", msg...)

	if rec.level&LInfo != 0 {
		log := GenLog(logs)
		green := newColor(FgBrGreen)
		printLog(green, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Warn(msg ...interface{}) {
	logs := genLog("Warn", msg...)

	if rec.level&LInfo != 0 {
		log := GenLog(logs)
		magenta := newColor(FgBrYellow)
		printLog(magenta, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Error(msg ...interface{}) {
	logs := genLog("Error", msg...)

	if rec.level&LError != 0 {
		log := GenLog(logs)
		red := newColor(FgBrRed)
		printLog(red, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Panic(msg ...interface{}) {
	logs := genLog("Panic", msg...)

	if rec.level&LError != 0 {
		log := GenLog(logs)
		red := newColor(FgBrMagenta)
		printLog(red, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
	panic(msg)
}

func (rec Logger) HTTP(msg ...interface{}) {
	logs := genLog("HTTP", msg...)

	if rec.level&LDebug != 0 {
		log := GenLog(logs)
		color := newColor(FgBrWhite)
		printLog(color, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Tracef(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Trace", s)

	if rec.level&LTrace != 0 {
		log := GenLog(logs)
		color := newColor(FgBrBlue)
		printLog(color, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Debugf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Debug", s)

	if rec.level&LDebug != 0 {
		log := GenLog(logs)
		color := newColor(FgBrCyan)
		printLog(color, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Infof(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Info", s)

	if rec.level&LInfo != 0 {
		log := GenLog(logs)
		green := newColor(FgBrGreen)
		printLog(green, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Warnf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Warn", s)

	if rec.level&LInfo != 0 {
		log := GenLog(logs)
		magenta := newColor(FgBrYellow)
		printLog(magenta, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Errorf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Error", s)

	if rec.level&LError != 0 {
		log := GenLog(logs)
		red := newColor(FgBrRed)
		printLog(red, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}

func (rec Logger) Fatal(msg ...interface{}) {
	msgs := append(msg, "os.Exit(1)")
	logs := genLog("Fatal", msgs...)

	if rec.level&LInfo != 0 {
		log := GenLog(logs)
		yellow := newColor(FgRed, BgBrBlue)
		printLog(yellow, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
	os.Exit(1)
}

func (rec Logger) Fatalf(format string, msg ...interface{}) {
	msgs := append(msg, "os.Exit(1)")
	s := fmt.Sprintf(format, msgs...)
	logs := genLog("Fatal", s)

	if rec.level&LInfo != 0 {
		log := GenLog(logs)
		yellow := newColor(FgRed, BgBrBlue)
		printLog(yellow, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
	os.Exit(1)
}

func (rec Logger) HTTPf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("HTTP", s)

	if rec.level&LHTTP != 0 {
		log := GenLog(logs)
		color := newColor(FgBrWhite)
		printLog(color, log)
		if rec.save {
			writeToFile(rec.file, log)
		}
	}
}
