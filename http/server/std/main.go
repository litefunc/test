package main

import (
	"log"
	"net/http"
	"test/http/server/std/internal"
)

func main() {

	http.HandleFunc("/formfile", internal.FormFile)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
