package main

import (
	"cloud/lib/logger"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path"
	"strings"
	"test/redmine-gitlab/gitlab"
	"time"
)

type HTTP struct {
	Token string   `json:"token"`
	Repos []string `json:"repos"`
}

type Git struct {
	Repos []string `json:"repos"`
}

type Config struct {
	LogLevel int  `json:"logger_level"`
	Tick     int  `json:"tick"`
	HTTP     HTTP `json:"http"`
	Git      Git  `json:"git"`
}

func main() {

	cfgFile := flag.String("cfg", "/root/config.json", "")
	repoDir := flag.String("data", "/root/data/redmine/repos", "")

	by, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		logger.Error(err)
		return
	}
	var cfg Config
	if err := json.Unmarshal(by, &cfg); err != nil {
		logger.Error(err)
		return
	}
	logger.SetLevel(cfg.LogLevel)

	cli := gitlab.NewClient(cfg.HTTP.Token)

	seconds := time.Second * time.Duration(cfg.Tick)
	ticker := time.NewTicker(seconds)

	for _, v := range cfg.HTTP.Repos {
		cli.CloneWithHTTP(*repoDir, v)
	}
	for _, v := range cfg.Git.Repos {
		gitlab.Clone(*repoDir, v)
	}

	for _ = range ticker.C {

		for _, v := range cfg.HTTP.Repos {

			if err := cli.CloneWithHTTP(*repoDir, v); err != nil {
				continue
			}

			ss := strings.Split(v, "/")
			if n := len(ss); n != 0 {

				dir := path.Join(*repoDir, ss[n-1])
				gitlab.Fetch(dir)

			}
		}

		for _, v := range cfg.Git.Repos {

			if err := gitlab.Clone(*repoDir, v); err != nil {
				continue
			}

			ss := strings.Split(v, "/")
			if n := len(ss); n != 0 {

				dir := path.Join(*repoDir, ss[n-1])
				gitlab.Fetch(dir)

			}
		}

	}

}
