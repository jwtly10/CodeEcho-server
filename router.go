package main

import "net/http"

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

type Routes []Route

func NewRouter(m *Middleware, h *Handlers) *http.ServeMux {
	router := http.NewServeMux()

	routes := Routes{
		Route{Path: "/api/v1/chatgpt/stream", Handler: m.HandleMiddleware(h.ChatGPTStreamHandler)},
		Route{Path: "/api/v1/transcribe", Handler: m.HandleMiddleware(h.DeepGramTranscribeHandler)},
	}

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler)
	}

	return router
}
