package main

import (
	"encoding/json"
	"fmt"
)

type File struct {
	Filesize uint64 `json:"filesize"`
	Filename string `json:"filename"`
}

type file struct {
	filesize uint64 `json:"filesize"`
	filename string `json:"filename"`
}

func public() {

	f := File{
		Filesize: 123,
		Filename: "http://192.168.2.2:40000/vod/noovo/1775-%E6%B5%B7%E7%B6%BF%E5%AF%B6%E5%AF%B6%EF%BC%9A%E5%A5%94%E8%B7%91%E5%90%A7/%E6%B5%B7%E7%B6%BF%E5%AF%B6%E5%AF%B6%EF%BC%9A%E5%A5%94%E8%B7%91%E5%90%A7.png",
	}

	by, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(by))

	var o File
	json.Unmarshal(by, &o)
	fmt.Println(o)

}

func private() {

	f := file{
		filesize: 123,
		filename: "http://192.168.2.2:40000/vod/noovo/1775-%E6%B5%B7%E7%B6%BF%E5%AF%B6%E5%AF%B6%EF%BC%9A%E5%A5%94%E8%B7%91%E5%90%A7/%E6%B5%B7%E7%B6%BF%E5%AF%B6%E5%AF%B6%EF%BC%9A%E5%A5%94%E8%B7%91%E5%90%A7.png",
	}

	by, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(by))

	var o file
	json.Unmarshal(by, &o)
	fmt.Println(o)

	s := `{"filesize":123,"filename":"http://192.168.2.2:40000/vod/noovo/1775-%E6%B5%B7%E7%B6%BF%E5%AF%B6%E5%AF%B6%EF%BC%9A%E5%A5%94%E8%B7%91%E5%90%A7/%E6%B5%B7%E7%B6%BF%E5%AF%B6%E5%AF%B6%EF%BC%9A%E5%A5%94%E8%B7%91%E5%90%A7.png"}`
	fmt.Println(s)

	json.Unmarshal([]byte(s), &o)
	fmt.Println(o)

	var o1 File
	json.Unmarshal([]byte(s), &o1)
	fmt.Println(o)
}

func main() {

	fmt.Println("public")
	public()

	fmt.Println("private")
	private()
}
