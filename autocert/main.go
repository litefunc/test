package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// setup a simple handler which sends a HTHS header for six months (!)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Secure World")
	})

	// look for the domains to be served from command line args
	flag.Parse()
	domains := flag.Args()
	for i := range domains {
		log.Println("domain", i+1, domains[i])
	}

	// create the autocert.Manager with domains and path to the cache
	var certManager autocert.Manager
	if len(domains) == 0 {
		log.Println(0)
		certManager = autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("/var/lib/.cache/golang-autocert"),
		}
	} else {
		certManager = autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      autocert.DirCache("/var/lib/.cache/golang-autocert"),
		}
	}

	// create the server itself
	server := &http.Server{
		Addr:    ":https",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Printf("Serving http/https for domains: %+v", domains)
	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		h := certManager.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	// serve HTTPS!
	log.Fatal(server.ListenAndServeTLS("", ""))

}
