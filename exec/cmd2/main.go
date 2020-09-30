package main

import (
	"bufio"
	"mstore/logger"
	"os/exec"
	"path"
)

func cmd(i int) *exec.Cmd {

	// p := path.Join(os.Getenv("GOPATH"), "src/test/exec/cmd1/cmd1")

	p := path.Join("/home/david/program/go/src/test/exec/cmd1/cmd1")

	cmd := exec.Command("sh", "-c", p)
	logger.Debug(cmd.String())
	rc, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error(err)
		return nil
	}
	if err := cmd.Start(); err != nil {
		logger.Error(err)
		return nil
	}

	scanner := bufio.NewScanner(rc)
	go func() {
		for scanner.Scan() {
			logger.Debug(i, scanner.Text())
		}
	}()

	return cmd
}

func sh(args ...string) {
	arg := append([]string{"-c"}, args...)
	cmd := exec.Command("sh", arg...)
	logger.Debug(cmd.String())
	by, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err, string(by))
		return
	}
	logger.Debug(string(by))
}

func Exec(name string, arg ...string) {

	cmd := exec.Command(name, arg...)
	logger.Debug(cmd.String())
	by, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err, string(by))
		return
	}
	logger.Debug(string(by))
}

func main() {
	sh("echo 1")
	sh("echo", "1")
	sh("'echo 1'")
	sh(`"echo 1"`)
	sh("echo '")
	sh("ls | grep main.go")
	Exec("ls", "| grep main.go")

	cmd1 := cmd(1)
	// cmd2 := cmd(2)

	cmd1.Wait()
	// cmd2.Wait()
}
