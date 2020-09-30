package main

import (
	"MediaImage/logger"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Compose struct {
	Version  string           `yaml:"version"`
	Services map[string]Image `yaml:"services"`
}

type Image struct {
	Image      string   `yaml:"image"`
	Ports      []string `yaml:"ports"`
	Volumes    []string `yaml:"volumes"`
	Privileged bool     `yaml:"privileged"`
}

func main() {

	p := path.Join(os.Getenv("GOPATH"), "src/test/yaml/docker-compose.yml")
	by, err := ioutil.ReadFile(p)
	if err != nil {
		logger.Error(err)
		return
	}

	var t Compose

	if err := yaml.Unmarshal(by, &t); err != nil {
		logger.Error(err)
		return
	}
	by1, err := json.Marshal(t)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by1))

	for i, v := range t.Services {
		logger.Debug(i, v)
	}

}
