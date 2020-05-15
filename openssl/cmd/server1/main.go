package main

import (
	"flag"
	"fmt"
	"log"
	"mstore/logger"
	"net/http"
	"os"
	"path"
	"strconv"
)

func main() {

	p := flag.Int("p", 80, "port")

	gopath := os.Getenv("GOPATH")
	cert := flag.String("cert", path.Join(gopath, "/src/test/openssl/static/server.crt"), "cert")
	key := flag.String("key", path.Join(gopath, "/src/test/openssl/static/server.key"), "key")

	flag.Parse()
	port := strconv.Itoa(*p)

	logger.Debug(port)
	logger.Debug(*cert)
	logger.Debug(*key)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	http.Handle("/static/", http.FileServer(http.Dir("./static")))

	go http.ListenAndServe(":http", httpToHTTPS())
	log.Fatal("HTTPS server error: ", http.ListenAndServeTLS(":"+port, *cert, *key, nil))

	// if *p == 443 {
	// 	go http.ListenAndServe(":http", httpToHTTPS())
	// 	log.Fatal("HTTPS server error: ", http.ListenAndServeTLS(":"+port, *cert, *key, nil))

	// } else {
	// 	log.Fatal("HTTP server error: ", http.ListenAndServe(":"+port, nil))
	// }
}

// func main() {

// 	p := flag.Int("p", 80, "port")

// 	gopath := os.Getenv("GOPATH")
// 	cert := flag.String("cert", path.Join(gopath, "/src/test/openssl/static/server.crt"), "cert")
// 	key := flag.String("key", path.Join(gopath, "/src/test/openssl/static/server.key"), "key")

// 	flag.Parse()

// 	port := strconv.Itoa(*p)

// 	logger.Debug(*cert)
// 	logger.Debug(*key)

// 	router := gin.Default()

// 	go router.RunTLS(":"+port, *cert, *key)
// 	router.Run(":80")

// }

// redirect connection form http to https
func httpToHTTPS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusFound)
	}
}
