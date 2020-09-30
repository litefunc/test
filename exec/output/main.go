package main

import (
	"bytes"
	"mstore/logger"
	"os"
	"os/exec"
	"path"
)

func main() {

	p := path.Join(os.Getenv("GOPATH"), "src/test/exec/cmd3/cmd3")
	cmd := exec.Command(p)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	by, err := cmd.Output()
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(by))
	logger.Error(stderr.String())

	// time.Sleep(time.Second * 5)
}
