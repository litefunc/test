package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	var p string
	flag.StringVar(&p, "p", ":9000", "port")
	flag.Parse()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// fs1 := http.FileServer(http.Dir("./static1"))
	// http.Handle("/static1/", fs1)

	err := http.ListenAndServe(p, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
