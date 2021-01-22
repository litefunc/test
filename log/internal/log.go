package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
)

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

	// for arm
	// err := syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)

	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Printf("Failed to redirect stderr to file: %v", err)
	}
}

type Logger struct {
	*log.Logger
}

func NewLogger(l *log.Logger) *Logger {
	return &Logger{l}
}

func (rec *Logger) Println(v ...interface{}) {
	rec.Logger.Println(v...)
}

func (rec *Logger) Debug(v ...interface{}) {
	rec.Logger.SetPrefix("DEBUG ")
	ss := make([]string, len(v), len(v))
	for i := range v {
		ss[i] = fmt.Sprintf(`%v`, v[i])
	}
	rec.Logger.Output(2, strings.Join(ss, " "))
}

func (rec *Logger) Info(v ...interface{}) {
	rec.Logger.SetPrefix("INFO  ")
	ss := make([]string, len(v), len(v))
	for i := range v {
		ss[i] = fmt.Sprintf(`%v`, v[i])
	}
	rec.Logger.Output(2, strings.Join(ss, " "))
}
