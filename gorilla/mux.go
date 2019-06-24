package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]

	fmt.Fprintf(w, "hello, %s!\n", key)
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]

	fmt.Fprintf(w, "hello, %s!\n", category)
}

func ArticleCategoryKeyHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]
	key := vars["key"]

	fmt.Fprintf(w, "hello, %s! %s\n", category, key)
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]
	id := vars["id"]

	fmt.Fprintf(w, "hello, %s! %s\n", category, id)
}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", Index)
	r.HandleFunc("/products/{key}", ProductHandler)
	r.HandleFunc("/articles/{category}", ArticlesCategoryHandler)
	r.HandleFunc("/articles/{category}/{key}", ArticleCategoryKeyHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	r.HandleFunc("/articles/test/test/test", Index)
	http.Handle("/", r)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/home/david/Downloads/noovo/MSAT/storage/"))))

	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// 	http.ServeFile(w, r, "index.html")
	// })

	// log.Fatal(http.ListenAndServe(":8080", r))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
