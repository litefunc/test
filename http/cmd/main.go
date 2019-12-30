package main

import (
	"cloud/lib/logger"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type hd struct {
	i  int
	mu sync.Mutex
}

func (h *hd) Get(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	logger.Debug(h.i)
	h.i = h.i + 1
	h.mu.Unlock()
	time.Sleep(time.Second)
	h.Sleep()
}

func (h *hd) Sleep() {
	time.Sleep(time.Second)
}

func main() {

	p := flag.Int("p", 8080, "port")
	flag.Parse()

	var n int
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n = n + 1
		hello(w, r)
		fmt.Fprintf(w, "visits: %d\n", n)
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		n = n + 1
		hello(w, r)
		fmt.Fprintf(w, "visits: %d\n", n)
	})

	http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get(fmt.Sprintf("http://localhost:%d/index", *p))
		if err != nil {
			return
		}
		by, _ := ioutil.ReadAll(res.Body)
		w.Write(by)
	})

	http.HandleFunc("/go1", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/index", *p), nil)
		if err != nil {
			return
		}
		req.Header.Add("X-Real-IP", r.RemoteAddr)
		res, err := client.Do(req)
		if err != nil {
			return
		}

		by, _ := ioutil.ReadAll(res.Body)
		w.Write(by)
	})

	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	})

	var h hd
	http.HandleFunc("/sleep", h.Get)

	var mu sync.Mutex
	var s int
	http.HandleFunc("/sleep1", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		logger.Debug(s)
		s++
		mu.Unlock()
		time.Sleep(time.Second)
	})

	http.HandleFunc("/sleep100", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sleep 100"))
		time.Sleep(time.Second * 100)
	})

	fmt.Println("HTTP server listen at:", *p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%d`, *p), nil))

}

func hello(w http.ResponseWriter, r *http.Request) {

	var names []string
	m := make(map[string][]string)
	for name, values := range r.Header {
		names = append(names, name)
		m[name] = values
	}
	sort.Strings(names)
	for _, name := range names {
		values := m[name]
		// Loop over all values for the name.
		for _, value := range values {

			fmt.Fprintf(w, "%s: %s\n\n", name, value)
		}
	}

	fmt.Fprintf(w, "r.Method : %s\n", r.Method)
	fmt.Fprintf(w, "r.Proto: %s\n", r.Proto)
	fmt.Fprintf(w, "r.Host: %s\n", r.Host)
	fmt.Fprintf(w, "r.URL.Host: %s\n", r.URL.Host)
	fmt.Fprintf(w, "r.URL.Path: %q\n", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "r.RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "r.RequestURI : %s\n", r.RequestURI)
}
