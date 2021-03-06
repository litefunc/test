package main

import (
	"cloud/lib/logger"
	"fmt"
	"os"
	"test/template/cloud/ota/service"
	model "test/template/model/xorm"
	"text/template"
)

func gen(text string, f *os.File, data interface{}) {

	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		panic(err)
	}

}

func stdout(text string, ser interface{}) {

	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, ser)
	if err != nil {
		panic(err)
	}

}

func genModel(tmp string, md model.Model) {

	dir := os.Getenv("GOPATH") + "/src/test/template"

	sdir := dir + "/out/model"
	if err := os.MkdirAll(sdir, os.ModePerm); err != nil {
		logger.Panic(err)
	}

	f, err := os.OpenFile(fmt.Sprintf(`%s/%s`, sdir, md.File), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gen(tmp, f, md)
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
		// service.New("msfw", "MsFw", "MsFws"),
		// service.New("routerfw", "RouterFw", "RouterFws"),
		// service.New("docker", "Docker", "Dockers"),
		// service.New("notification", "Notification", "Notifications"),
		// service.New("ui_share_com", "UIShareCom", "UIShareComs"),
		// service.New("ui_share_group", "UIShareGroup", "UIShareGroups"),
		// service.New("system", "System", "Systems"),
		// service.New("group_system", "GroupSystem", "GroupSystems"),
		// service.New("system_share_com", "SystemShareCom", "SystemShareComs"),
		// service.New("system_share_group", "SystemShareGroup", "SystemShareGroups"),
	}

	inits(sers)
	get(sers)
	create(sers)
	update(sers)
	delete(sers)

	genModel(model.Tmpl, model.Md)
}
