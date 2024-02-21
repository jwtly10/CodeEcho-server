package main

import "net/http"

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func NewRouter(handlers *Handlers) *http.ServeMux {
	router := http.NewServeMux()

	routes := Routes{
		Route{"/", handlers.HomeHandler},
	}

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler)
	}

	return router
}
