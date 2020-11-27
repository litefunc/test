package main

import (
	"VodoPlay/logger"
	"bytes"

	"os/exec"
	"path"
)

func main() {

	// p := path.Join(os.Getenv("GOPATH"), "src/test/exec/cmd3/cmd3")
	// cmd := exec.Command(p)

	p := path.Join("/home/david/program/docker-compose/composetest/docker-compose.yml")
	// cmd := exec.Command("docker-compose", "-f", p, "up", "-d", "--remove-orphans")
	cmd := exec.Command("docker-compose", "-f", p, "ps")
	// cmd := exec.Command("docker-compose", "-f", p, "down")

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
