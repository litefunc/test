package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"test/logger"
)

func main() {

	u := "https://www.google.com/"

	// u := "http://noovo-dock.ddns.net/redmine/"

	remote, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	director := func(req *http.Request) {
		// req.Header.Add("X-Forwarded-Host", req.Host)
		// req.Header.Add("X-Origin-Host", remote.Host)
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		logger.Debug(req.URL)
	}

	proxy := &httputil.ReverseProxy{Director: director}

	// proxy := httputil.NewSingleHostReverseProxy(remote)
	// use http.Handle instead of http.HandleFunc when your struct implements http.Handler interface
	// http.Handle("/redmine/", &ProxyHandler{proxy})
	http.Handle("/", &ProxyHandler{proxy})
	http.HandleFunc("/123", index)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type ProxyHandler struct {
	p *httputil.ReverseProxy
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("X-Ben", "Rad")
	ph.p.ServeHTTP(w, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(1, r.URL)

}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-Ben", "Rad")
		p.ServeHTTP(w, r)
	}
}
