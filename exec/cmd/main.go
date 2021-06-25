package main

import (
	"LocalServer/logger"
	"bytes"
	"os/exec"
)

func main() {

	cmd := exec.Command("pkill", "-f", "main.go")

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	if err := cmd.Start(); err != nil {
		logger.Error(err)
		return
	}
	if err := cmd.Wait(); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(b.String())
}
