package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	i := flag.Int("i", 1, "server number")
	p := flag.Int("p", 8088, "port")
	flag.Parse()

	var n int
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n = n + 1
		fmt.Fprintf(w, "server: %d\n", *i)
		fmt.Fprintf(w, "r.Host: %s\n", r.Host)
		fmt.Fprintf(w, "r.URL.Host: %s\n", r.URL.Host)
		fmt.Fprintf(w, "r.URL.Path: %q\n", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "visits: %d\n", n)
	})

	fmt.Println("HTTP server listen at:", *p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%d`, *p), nil))

}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func tls() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
