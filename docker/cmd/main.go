package main

import (
	"cloud/lib/logger"
	"errors"
	"os/exec"
	"strings"
)

func cmd(args ...string) (string, error) {
	out, err := exec.Command("docker", args...).CombinedOutput()
	logger.Warn("docker", strings.Join(args, " "))
	if err != nil && len(out) != 0 {
		logger.Error(err.Error() + ": " + string(out))
		return "", errors.New(err.Error() + ": " + string(out))
	}

	return strings.TrimSuffix(string(out), "\n"), err
}

func prune() error {

	args := []string{"system", "prune", "-a", "-f"}

	// find containers which based on img
	out, err := cmd(args...)
	if err != nil {
		return err
	}
	logger.Debug(out)
	return nil
}

func main() {
	logger.Debug(prune())
}
