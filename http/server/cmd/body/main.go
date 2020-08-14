package main

import (
	"cloud/lib/logger"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	p := flag.Int("p", 8080, "port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hello(w, r)
	})

	fmt.Println("HTTP server listen at:", *p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%d`, *p), nil))

}

func hello(w http.ResponseWriter, r *http.Request) {

	by, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(by))
	w.Write(by)
}
