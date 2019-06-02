package router

import (
	"github.com/gorilla/mux"
)

type Route struct {
	*mux.Route
}

func NewRoute(r *mux.Route) *Route {
	return &Route{r}
}
