package internal

import (
	"log"
	"os"
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
