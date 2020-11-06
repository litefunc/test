package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	var p string
	flag.StringVar(&p, "p", ":8080", "port")
	flag.Parse()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	err := http.ListenAndServe(p, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
