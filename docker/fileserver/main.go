package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	var p string
	flag.StringVar(&p, "p", ":8087", "port")
	flag.Parse()

	fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", fs)

	// http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(p, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
