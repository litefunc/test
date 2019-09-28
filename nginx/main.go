package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"test/nginx/grpc/hello"
	"time"
)

func main() {
	i := flag.Int("i", 1, "server number")
	p := flag.Int("p", 8080, "port")
	hp := flag.Int("hp", 50050, "grpc hello server port")
	flag.Parse()

	hello.NewServer(uint64(*i), *hp)

	var n int

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n = n + 1
		fmt.Fprintf(w, "server: %d\n", *i)
		fmt.Fprintf(w, "r.Host: %s\n", r.Host)
		fmt.Fprintf(w, "r.URL.Host: %s\n", r.URL.Host)
		fmt.Fprintf(w, "r.URL.Path: %q\n", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "visits: %d\n", n)
	})

	http.HandleFunc("/wait", func(w http.ResponseWriter, r *http.Request) {
		n = n + 1
		fmt.Fprintf(w, "server: %d\n", *i)
		fmt.Fprintf(w, "r.Host: %s\n", r.Host)
		fmt.Fprintf(w, "r.URL.Host: %s\n", r.URL.Host)
		fmt.Fprintf(w, "r.URL.Path: %q\n", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "visits: %d\n", n)
		w.Write([]byte("wait\n"))
		go func() {
			time.Sleep(time.Second * 4)
			fmt.Fprintln(w, "finish 1")
		}()
		go func() {
			time.Sleep(time.Second * 6)
			fmt.Fprintln(w, "finish 2")
		}()
		time.Sleep(time.Second * 5)
		fmt.Fprintln(w, "return")
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
