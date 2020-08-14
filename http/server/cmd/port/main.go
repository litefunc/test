package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	p := flag.Int("p", 8090, "port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hello(w, r)
	})

	fmt.Println("HTTP server listen at:", *p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%d`, *p), nil))

}

func hello(w http.ResponseWriter, r *http.Request) {

	// var wc chan struct{}
	// <-wc
	w.Write([]byte("ok"))
}
