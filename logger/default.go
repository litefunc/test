package logger

import (
	"fmt"
	"os"
)

var defaultLogger = Logger{
	flag:  LTrace | Ltime | Lfile | Lline,
	level: LTrace | LDebug | LInfo | LWarn | LError | LFatal | LPanic | LHTTP,
	save:  false,
}

func GetDefaultLogger() *Logger {
	return &defaultLogger
}

func Trace(msg ...interface{}) {
	logs := genLog("Trace", msg...)

	if defaultLogger.level&LTrace != 0 {
		log := GenLog(logs)
		color := newColor(FgBrBlue)
		printLog(color, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Debug(msg ...interface{}) {
	logs := genLog("Debug", msg...)

	if defaultLogger.level&LDebug != 0 {
		log := GenLog(logs)
		color := newColor(FgBrCyan)
		printLog(color, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Info(msg ...interface{}) {
	logs := genLog("Info", msg...)

	if defaultLogger.level&LInfo != 0 {
		log := GenLog(logs)
		green := newColor(FgBrGreen)
		printLog(green, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Warn(msg ...interface{}) {
	logs := genLog("Warn", msg...)

	if defaultLogger.level&LInfo != 0 {
		log := GenLog(logs)
		magenta := newColor(FgBrYellow)
		printLog(magenta, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Error(msg ...interface{}) {
	logs := genLog("Error", msg...)

	if defaultLogger.level&LError != 0 {
		log := GenLog(logs)
		red := newColor(FgBrRed)
		printLog(red, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Panic(msg ...interface{}) {
	logs := genLog("Panic", msg...)

	if defaultLogger.level&LError != 0 {
		log := GenLog(logs)
		red := newColor(FgBrMagenta)
		printLog(red, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
	panic(msg)
}

func HTTP(msg ...interface{}) {
	logs := genLog("HTTP", msg...)

	if defaultLogger.level&LDebug != 0 {
		log := GenLog(logs)
		color := newColor(FgBrWhite)
		printLog(color, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Tracef(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Trace", s)

	if defaultLogger.level&LTrace != 0 {
		log := GenLog(logs)
		color := newColor(FgBrBlue)
		printLog(color, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Debugf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Debug", s)

	if defaultLogger.level&LDebug != 0 {
		log := GenLog(logs)
		color := newColor(FgBrCyan)
		printLog(color, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Infof(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Info", s)

	if defaultLogger.level&LInfo != 0 {
		log := GenLog(logs)
		green := newColor(FgBrGreen)
		printLog(green, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Warnf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Warn", s)

	if defaultLogger.level&LInfo != 0 {
		log := GenLog(logs)
		magenta := newColor(FgBrYellow)
		printLog(magenta, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Errorf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("Error", s)

	if defaultLogger.level&LError != 0 {
		log := GenLog(logs)
		red := newColor(FgBrRed)
		printLog(red, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Fatal(msg ...interface{}) {
	msgs := append(msg, "os.Exit(1)")
	logs := genLog("Fatal", msgs...)

	if defaultLogger.level&LInfo != 0 {
		log := GenLog(logs)
		yellow := newColor(FgRed, BgBrBlue)
		printLog(yellow, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
	os.Exit(1)
}

func Fatalf(format string, msg ...interface{}) {
	msgs := append(msg, "os.Exit(1)")
	s := fmt.Sprintf(format, msgs...)
	logs := genLog("Fatal", s)

	if defaultLogger.level&LInfo != 0 {
		log := GenLog(logs)
		yellow := newColor(FgRed, BgBrBlue)
		printLog(yellow, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
	os.Exit(1)
}

func HTTPf(format string, msg ...interface{}) {
	s := fmt.Sprintf(format, msg...)
	logs := genLog("HTTP", s)

	if defaultLogger.level&LHTTP != 0 {
		log := GenLog(logs)
		color := newColor(FgBrWhite)
		printLog(color, log)
		if defaultLogger.save {
			writeToFile(defaultLogger.file, log)
		}
	}
}

func Exit(err error) {
	if err != nil {
		msgs := []interface{}{err.Error()}
		logs := genLog("Exit", msgs...)

		if defaultLogger.level&LError != 0 {
			log := GenLog(logs)
			red := newColor(FgBrRed)
			printLog(red, log)
			if defaultLogger.save {
				writeToFile(defaultLogger.file, log)
			}
		}
		os.Exit(1)
	}
}

func LogErr(err error) {
	if err != nil {
		msgs := []interface{}{err.Error()}
		logs := genLog("Error", msgs...)

		if defaultLogger.level&LError != 0 {
			log := GenLog(logs)
			red := newColor(FgBrRed)
			printLog(red, log)
			if defaultLogger.save {
				writeToFile(defaultLogger.file, log)
			}
		}
	}
}
