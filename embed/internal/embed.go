package internal

import (
	"embed"
	"test/logger"
)

//go:embed hello.txt
var s string

//go:embed hello.txt
var b []byte

//go:embed hello.txt
var f embed.FS

//go:embed dir
var d embed.FS

func Embed() {

	logger.Debug(s)

	logger.Debug(string(b))

	data, err := f.ReadFile("hello.txt")
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(data))

	data, err = d.ReadFile("dir/a.txt")
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(data))

	fs, err := d.ReadDir("dir")
	if err != nil {
		logger.Error(err)
	}
	for _, v := range fs {

		logger.Debug(v.Name(), v.IsDir())
	}
}
