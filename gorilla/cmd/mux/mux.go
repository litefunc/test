package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)
 
type Handler func(http.ResponseWriter, *http.Request)

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r)
}

func simpleMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func simpleMw1(next Handler) Handler {

	f := func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	}
	return Handler(http.HandlerFunc(f))
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func test(i int) {
	fmt.Print(i)
}

func testf(i int, f func(int)) {
	f(i)
}

type tf func(int)

type I int

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/articles/{category}/", simpleMw1(ArticlesCategoryHandler))
	// r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

	// simpleMw(ArticlesCategoryHandler)

	r.Use(simpleMw)

	// log.Fatal(http.ListenAndServe(":8088", nil))

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	srv1 := &http.Server{
		Handler: r,
		Addr:    ":8088",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv1.ListenAndServe())

	var i tf

	testf(0, i)
}
