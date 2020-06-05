package main

import (
	"VodoPlay/logger"
	"net/url"
	"path"
)

const (
	p1 = "localhost"
	p2 = "localhost:80"
	p3 = "http://localhost"
	p4 = "http://localhost:8080"
	p5 = "192.168.2.2"
	p6 = "192.168.2.2:8080"
	p7 = "http://192.168.2.2"
	p8 = "http://192.168.2.2:8080"
	p9 = "http://192.168.2.2:8080"

	p10 = "http:/192.168.2.2:8080"
	p11 = "http:///192.168.2.2:8080"

	p12 = "http://192.168.2.2:40000/media/abc  def  gh.mp4"
	p13 = "http://192.168.2.2:40000/media/abc%20%20def%20%20gh.mp4"
)

func main() {

	// prefix := "http://localhost:8080"
	// u, err := url.Parse(prefix)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// u1 := *u
	// u1.Path = path.Join(u1.Path, "abc")

	// u2 := *u
	// u2.Path = path.Join(u2.Path, "def")

	// logger.Debug(u1.String())
	// logger.Debug(u2.String())

	us(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10)
	us1(p1, p2, p3, p4, p5, p6, p7, p8, p9, p11)
	us2(p1, p2, p3, p4, p5, p6, p7, p8, p9, p11)
	us3(p12, p13, "http://192.168.2.2:40000/file/noovo/files/布魯克林孤兒 Motherless Brooklyn.json")

	u, err := url.Parse("http://192.168.2.2:40000")
	if err != nil {
		logger.Error(err)
	}

	u.Path = path.Join(u.Path, "file/noovo/files/布魯克林孤兒 Motherless Brooklyn.json")
	us3(u.String())
}

func us(ps ...string) {
	for i, p := range ps {
		logger.Debug(i, p)
		u, err := url.Parse(p)
		if err != nil {
			logger.Error(err)
			logger.Error(u == nil)
			continue
		}

		u.Path = path.Join(u.Path, "abc")
		logger.Debug(u.String())

	}
}

func us1(ps ...string) {
	for i, p := range ps {
		logger.Debug(i, p)
		u, err := url.Parse(p)
		if err != nil {
			logger.Error(err)
			logger.Error(u == nil)
			continue
		}

		u.Path = path.Join(u.Path, "//abc//")
		logger.Debug(u.String())

	}
}

func us2(ps ...string) {
	for i, p := range ps {
		logger.Debug(i, p)
		u, err := url.Parse(p)
		if err != nil {
			logger.Error(err)
			logger.Error(u == nil)
			continue
		}

		u.Path = path.Join(u.Path, " a b c 中文")
		logger.Debug(u.String())

	}
}

func us3(ps ...string) {
	for i, p := range ps {
		logger.Debug(i, p)
		u, err := url.Parse(p)
		if err != nil {
			logger.Error(err)
			logger.Error(u == nil)
			continue
		}
		logger.Debug(u.String())
		logger.Debug(u.String() == p)
	}
}
