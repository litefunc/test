package main

import (
	"os"
	"syscall"
	"test/logger"
)

func main() {
	p, err := os.FindProcess(0)
	if err != nil {
		logger.Error(err)
		return

	} else {

		if err := p.Signal(syscall.Signal(0)); err != nil {
			logger.Error(err)
			return
		}

	}

	logger.Debug(p.Pid)
	// ps, err := p.Wait()
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// logger.Debug(ps.String())
}
