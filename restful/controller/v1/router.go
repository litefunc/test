package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	r := mux.NewRouter()

	var h MethodNotAllowedHandler
	r.MethodNotAllowedHandler = h

	return &Router{Router: r}
}

type Route struct {
	*mux.Route
}

func NewRoute(r *mux.Route) *Route {
	return &Route{r}
}

type MethodNotAllowedHandler func(w http.ResponseWriter, req *http.Request)

func (MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Content-Type", "text/plain")

	code := http.StatusMethodNotAllowed
	msg := http.StatusText(code)

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(msg))
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func addMiddlewares(h http.HandlerFunc, ms []Middleware) http.HandlerFunc {
	for _, m := range ms {
		h = m(h)
	}

	return h
}

func (r Router) HandleFunc(path string, h http.HandlerFunc, ms ...Middleware) *Route {
	h = addMiddlewares(h, ms)
	route := r.Router.HandleFunc(path, h)
	return NewRoute(route)
}

func (r Router) GET(path string, h http.HandlerFunc, ms ...Middleware) *Route {
	route := r.HandleFunc(path, h, ms...).Methods("GET")
	return NewRoute(route)
}

func (r Router) POST(path string, h http.HandlerFunc, ms ...Middleware) *Route {
	route := r.HandleFunc(path, h, ms...).Methods("POST")
	return NewRoute(route)
}

func (r Router) PUT(path string, h http.HandlerFunc, ms ...Middleware) *Route {
	route := r.HandleFunc(path, h, ms...).Methods("PUT")
	return NewRoute(route)
}

func (r Router) DELETE(path string, h http.HandlerFunc, ms ...Middleware) *Route {
	route := r.HandleFunc(path, h, ms...).Methods("DELETE")
	return NewRoute(route)
}

func (r Router) OPTIONS(path string, h http.HandlerFunc, ms ...Middleware) *Route {
	route := r.HandleFunc(path, h, ms...).Methods("OPTIONS")
	return NewRoute(route)
}
