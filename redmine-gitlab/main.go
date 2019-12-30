package main

import (
	"bytes"
	"cloud/lib/logger"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type Repos struct {
	Repos []string `json:"repos"`
}

func main() {

	cfgFile := flag.String("cfg", "/root/config.json", "")
	repoDir := flag.String("data", "/root/data/redmine/repos", "")
	cycle := flag.Int("cycle", 10, "")

	by, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		logger.Error(err)
		return
	}
	var repos Repos
	if err := json.Unmarshal(by, &repos); err != nil {
		logger.Error(err)
		return
	}

	seconds := time.Second * time.Duration(*cycle)
	ticker := time.NewTicker(seconds)

	for _, v := range repos.Repos {
		clone(*repoDir, v)
	}

	for _ = range ticker.C {

		for _, v := range repos.Repos {

			if err := clone(*repoDir, v); err != nil {
				continue
			}

			ss := strings.Split(v, "/")
			if n := len(ss); n != 0 {

				dir := path.Join(*repoDir, ss[n-1])
				fetch(dir)

			}
		}

	}

}

func clone(dir, url string) error {

	if err := os.Chdir(dir); err != nil {
		return err
	}

	ss := strings.Split(url, "/")
	if n := len(ss); n != 0 {

		_, err := os.Stat(ss[n-1])
		if err != nil {
			if os.IsNotExist(err) {

				_, err := output("git", "clone", "--mirror", url)
				if err != nil {
					return err
				}
				return nil

			}

			logger.Error(err)
			return err
		}

		return nil
	}

	return fmt.Errorf(`invalid url:%s`, url)

}

func fetch(dir string) error {

	if err := os.Chdir(dir); err != nil {
		return err
	}

	_, err := output("git", "fetch", "--all")
	if err != nil {
		return err
	}

	return nil

}

func output(name string, arg ...string) (string, error) {
	logger.Debug(name, strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		e := fmt.Errorf(`%s: %s`, fmt.Sprint(err), stderr.String())
		logger.Debug(e)
		return out.String(), e
	}

	logger.Debug(out.String())
	return out.String(), nil
}

func run(name string, arg ...string) error {
	logger.Debug(name, strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		e := fmt.Errorf(`%s: %s`, fmt.Sprint(err), stderr.String())
		logger.Debug(e)
		return e
	}

	return nil
}
