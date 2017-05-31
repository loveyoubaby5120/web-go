package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route specifies a route.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	PathPrefix  string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
	MatcherFunc mux.MatcherFunc
}

// NewRoute creates a new route.
func NewRoute(name string, method string, pattern string, handler http.HandlerFunc) Route {
	return Route{
		Name:        name,
		Method:      method,
		Pattern:     pattern,
		HandlerFunc: handler,
	}
}
