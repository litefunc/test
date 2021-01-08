package main

import (
	"bufio"
	"cloud/lib/logger"
	"os"
	"path/filepath"
)

func isDir(path string) (bool, error) {

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Error(err)
			return false, err
		}
		return false, nil
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

func createFileIfNotExist(path string) error {
	dir, err := isDir(path)
	if err != nil {
		return err
	}
	if dir {
		if err := os.RemoveAll(path); err != nil {
			logger.Error(err)
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		logger.Error(err)
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f.Close()
	return nil
}

func main() {

}

func stat(path ...string) {
	for _, p := range path {
		if _, err := os.Stat(p); err != nil {
			logger.Error(err)
		}
	}
}

func bufioWrite(path string, data []byte) error {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	n, err := w.Write(data)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("wrote %d bytes\n", n)
	if err := w.Flush(); err != nil {
		logger.Error(err)
		return err
	}

	if err := f.Chmod(os.ModePerm); err != nil {
		logger.Error(err)
	}
	if err := f.Sync(); err != nil {
		logger.Error(err)
	}
	if err := f.Close(); err != nil {
		logger.Error(err)
	}

	return nil
}
