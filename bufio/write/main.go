package main

import (
	"VodoPlay/logger"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		m[i] = i
	}
	by, _ := json.Marshal(m)
	p := "a.txt"
	writeLines([]string{string(by)}, p)
}

func writeLines(lines []string, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer file.Close()
	fi, _ := file.Stat()
	logger.Debug(fi.Size())

	w := bufio.NewWriter(file)
	for _, line := range lines {
		n, err := fmt.Fprintln(w, line)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Debug(n)
	}
	if err := w.Flush(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
