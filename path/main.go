package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
)

func main() {
	fmt.Println(path.Join("a", "b", "/c/"))
	fmt.Println(path.Join("/a", "b", "/c/"))
	fmt.Println(path.Join("http://localhost:8080", "b", "/c/"))

	u, err := url.Parse("//")
	if err != nil {
		log.Println(err)
		return
	}
	u.Path = path.Join(u.Path, "bar.html")
	s := u.String()
	log.Println(s)

}
