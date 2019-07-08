package messenger

import (
	"LocalServer/config"
	"LocalServer/docker"
	"LocalServer/logger"
	"LocalServer/messenger/deploy"
	"LocalServer/messenger/pb"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/Hunsin/beaver"
)

const (
	// PathDocker saves dockers received from cloud, they are supposed to be run by localserver
	PathDocker = "/docker.json"

	// OptsPath saves containers running on board currently
	OptsPath = "/docker_opts.json"

	demoDefaultPath        = "/assets/demo.tar"
	mediaServerDefaultPath = "/assets/media_server.tar"
)

type DockerClient struct {
	Docker    docker.Client
	dockerHub config.DockerHub
	deploy    *deploy.Client
	optsPath  string
	listPath  string
}

func NewDockerClient(hub config.DockerHub, deploy *deploy.Client, listPath, optsPath string) *DockerClient {
	return &DockerClient{dockerHub: hub, deploy: deploy, listPath: listPath, optsPath: optsPath}
}

func (cli DockerClient) RunDefaultDockers() {
	cs, err := cli.Docker.ReadOpts(config.CfgDir+ OptsPath)
	if err !=nil {
		return
	}
	if cli.Docker.ServicesContains(cs, "demo") {
		 cli.RunDemoDefaultIfNoDemo()
	}
	if cli.Docker.ServicesContains(cs, "media_server") {
		 cli.RunMediaServerDefaultIfNoMediaServer()
	}
}

func (cli DockerClient) RunDemoDefaultIfNoDemo() error {
	return cli.Docker.RunDemoDefaultIfNoDemo(config.CfgDir + demoDefaultPath)
}
func (cli DockerClient) RunMediaServerDefaultIfNoMediaServer() error {
	return cli.Docker.RunMediaServerDefaultIfNoMediaServer(config.CfgDir + mediaServerDefaultPath)
}

func (dc DockerClient) dockerCreate(hub config.DockerHub, s docker.Service) (string, error) {
	// login
	err := dc.Docker.Login(hub)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	// create container
	id, err := dc.Docker.Create(s.Repo, s.Tag, s.Options...)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return id, nil
}

func (dc DockerClient) UpdateDockersFromFile() error {

	by, err := ioutil.ReadFile(dc.listPath)
	if err != nil {
		beaver.Error(err.Error())
		return err
	}
	var update pb.Update
	if err := json.Unmarshal(by, &update); err != nil {
		beaver.Error(err.Error())
		return err
	}

	var ns []docker.Service
	if err := json.Unmarshal(update.Data, &ns); err != nil {
		logger.Error(err)
		return err
	}

	return dc.updateDockers(ns)
}

func (dc DockerClient) UpdateDockers(update pb.Update) error {

	var ns []docker.Service
	if err := json.Unmarshal(update.Data, &ns); err != nil {
		logger.Error(err)
		return err
	}

	logger.Warn(ns)

	if err := beaver.JSON(&update).WriteFile(dc.listPath); err != nil {
		return err
	}

	return dc.updateDockers(ns)
}

func (dc DockerClient) updateDockers(ns []docker.Service) error {

	defer dc.RunDefaultDockers()

	if len(ns) != 0 {
		if err := dc.Docker.Login(dc.dockerHub); err != nil {
			logger.Error(err)
			return err
		}
	}

	// e is a list of failed installed docker Service
	e := []docker.Service{}

	for i := range ns {
		if err := dc.Docker.Pull(ns[i].Repo, ns[i].Tag); err != nil {
			logger.Error(err)
			ns[i].Error = err.Error()
			e = append(e, ns[i])
		}
	}

	// report error if pull docker fail and then return
	if len(e) != 0 {

		rp := pb.Report{To: "ota", Sn: dc.deploy.Account.Sn}
		rperr := Err{Type: "docker"}

		defaultErrors.Docker = e
		beaver.JSON(&defaultErrors).WriteFile(pathErrors)

		var errs []DeviceDockerErr
		for i := range e {
			errs = append(errs, DeviceDockerErr{Repo: e[i].Repo, Tag: e[i].Tag, Error: e[i].Error})
		}
		dat, err := json.Marshal(errs)
		if err != nil {
			logger.Error(err)
			return err
		}
		rperr.Data = dat

		dat, err = json.Marshal(rperr)
		if err != nil {
			logger.Error(err)
			return err
		}
		rp.Data = dat

		return dc.deploy.Error(&rp)

	}

	var containers []docker.Service

	// get all docker containers on board
	ss, err := dc.Docker.Services(dc.optsPath)
	if err != nil {
		return err
	}

	// for each Service in docker, find if exists in new list
	for i := range ss {
		keep := false
		opts := strings.Join(ss[i].Options, "")

		for j := range ns {
			if ss[i].Repo == ns[j].Repo &&
				ss[i].Tag == ns[j].Tag &&
				opts == strings.Join(ns[j].Options, "") {
				keep = true
				containers = append(containers, ss[i])
				break
			}
		}

		// needs remove
		if !keep {
			if err := dc.Docker.Remove(ss[i].Repo + ":" + ss[i].Tag); err != nil {
				containers = append(containers, ss[i])
			}
		}
	}

	logger.Warn(ns)
	for i := range ns {
		old := false
		opts := strings.Join(ns[i].Options, "")

		for j := range ss {
			if ss[j].Repo == ns[i].Repo &&
				ss[j].Tag == ns[i].Tag &&
				opts == strings.Join(ss[j].Options, "") {
				old = true
				break
			}
		}

		// new service
		if !old {
			id, err := dc.Docker.Create(ns[i].Repo, ns[i].Tag, ns[i].Options...)
			if err != nil {
				ns[i].Error = err.Error()
				e = append(e, ns[i])
				continue
			}
			containers = append(containers, ns[i])
			if err := dc.Docker.Start(id); err != nil {
				ns[i].Error = err.Error()
				e = append(e, ns[i])
			}
		}
	}
	logger.Info(containers)

	// save current containers to file
	docker.SaveOpts(containers, dc.optsPath)

	// report error
	rp := pb.Report{To: "ota", Sn: dc.deploy.Account.Sn}
	rperr := Err{Type: "docker"}
	if len(e) != 0 {
		defaultErrors.Docker = e
		beaver.JSON(&defaultErrors).WriteFile(pathErrors)

		var errs []DeviceDockerErr
		for i := range e {
			errs = append(errs, DeviceDockerErr{Repo: e[i].Repo, Tag: e[i].Tag, Error: e[i].Error})
		}
		dat, err := json.Marshal(errs)
		if err != nil {
			logger.Error(err)
			return err
		}
		rperr.Data = dat

	}
	dat, err := json.Marshal(rperr)
	if err != nil {
		logger.Error(err)
		return err
	}
	rp.Data = dat

	return dc.deploy.Error(&rp)
}

type Err struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DeviceDockerErr struct {
	Repo  string `json:"repo"`
	Tag   string `json:"tag"`
	Error string `json:"error"`
}
