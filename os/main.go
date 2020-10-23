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

	// pj := path.Join

	// p := os.Getenv("GOPATH") + "/src/test/os/testdir/docker.json"

	// dir, err := filepath.Abs(filepath.Dir(p))
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }

	// // if err := os.RemoveAll(dir); err != nil {
	// // 	logger.Error(err)
	// // }
	// if err := os.MkdirAll(dir, os.ModePerm); err != nil {
	// 	logger.Error(err)
	// }
	// f, err := os.OpenFile(p, os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	logger.Error(err)
	// }
	// defer f.Close()

	// stat(p, "/usr/", "usr/abc", "Behind%20My%20Life/Behind%20My%20Life.jpg")
	// p = os.Getenv("GOPATH") + "/src/test/os/" + "Behind%20My%20Life/Behind%20My%20Life.jpg"
	// stat(p)
	// stat(strings.Replace(p, "%20", " ", -1))

	// logger.Debug(isDir(dir))
	// logger.Debug(isDir(pj(dir, "docker.json")))

	// createFileIfNotExist(pj(dir, "test.json"))

	p1 := os.Getenv("GOPATH") + "/src/test/os/docker1.json"
	// f, err = os.OpenFile(p1, os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }

	// ioutil.WriteFile(p1, []byte(`123`), os.ModePerm)
	// f.Close()

	// f1, err := os.OpenFile(p1, os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// defer f1.Close()

	bufioWrite(p1, []byte(`abc`))
	// bufioWrite(p1, []byte(`789`))
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
