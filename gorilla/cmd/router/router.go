package main

import (
	"fmt"
	"log"
	"net/http"
	"test/gorilla/mux/router"
	"time"

	"github.com/gorilla/mux"
)

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func simpleMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func simpleMw1(next http.HandlerFunc) http.HandlerFunc {

	f := func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next(w, r)
	}
	return f
}

func main() {
	r := router.New()

	r.GET("/articles/{category}/", ArticlesCategoryHandler, simpleMw1, simpleMw1)
	r.POST("/articles/{category}/", ArticlesCategoryHandler, simpleMw1, simpleMw1)

	// simpleMw(ArticlesCategoryHandler)
	r.Use(simpleMw)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
