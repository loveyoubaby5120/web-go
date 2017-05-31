package httpserv

import (
	"jdy/pkg/api"
	"jdy/pkg/common"

	"github.com/gorilla/mux"
)

func addRoute(router *mux.Router, route api.Route) {
	r := router.
		Methods(common.SplitAndTrim(route.Method, ",")...).
		Name(route.Name).
		Handler(route.HandlerFunc)
	if route.Pattern != "" {
		r.Path(route.Pattern)
	}
	if route.PathPrefix != "" {
		r.PathPrefix(route.PathPrefix)
	}
	if route.MatcherFunc != nil {
		r.MatcherFunc(route.MatcherFunc)
	}
	if route.Handler != nil {
		r.Handler(route.Handler)
	}
}

// AddRoutes adds multiple routes to the router.
func AddRoutes(router *mux.Router, routes []api.Route) {
	for _, route := range routes {
		addRoute(router, route)
	}
}
