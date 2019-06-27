package main

import (
	"cloud/lib/logger"
	"fmt"
	"os"
	"test/template/cloud/ota/service"
	"text/template"
)

func gen(text string, f *os.File, ser service.Service) {

	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, ser)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f, ser)
	if err != nil {
		panic(err)
	}

}

func genService(text, out string, ser service.Service) {

	dir := os.Getenv("GOPATH") + "/src/test/template"

	sdir := dir + "/out/service"
	if _, err := os.Stat(sdir); os.IsNotExist(err) {
		if err := os.Mkdir(sdir, os.ModePerm); err != nil {
			logger.Panic(err)
		}
	}

	pdir := fmt.Sprintf(`%s/%s`, sdir, ser.Package)
	if _, err := os.Stat(pdir); os.IsNotExist(err) {
		if err := os.Mkdir(pdir, os.ModePerm); err != nil {
			logger.Panic(err)
		}
	}

	f, err := os.OpenFile(fmt.Sprintf(`%s/%s`, pdir, out), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gen(text, f, ser)
}

func inits(sers []service.Service) {
	for _, ser := range sers {
		genService(service.Init, "service.go", ser)
	}
}

func get(sers []service.Service) {
	for _, ser := range sers {
		genService(service.Get, "get.go", ser)
	}
}

func create(sers []service.Service) {
	for _, ser := range sers {
		genService(service.Create, "create.go", ser)
	}
}

func update(sers []service.Service) {
	for _, ser := range sers {
		genService(service.Update, "update.go", ser)
	}
}

func delete(sers []service.Service) {
	for _, ser := range sers {
		genService(service.Delete, "delete.go", ser)
	}
}

func main() {

	sers := []service.Service{
		service.New("msfw", "MsFw", "MsFws"),
		service.New("routerfw", "RouterFw", "RouterFws"),
		service.New("docker", "Docker", "Dockers"),
		service.New("notification", "Notification", "Notifications"),
	}

	inits(sers)
	get(sers)
	create(sers)
	update(sers)
	delete(sers)
}
