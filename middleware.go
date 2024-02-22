package main

import (
	"log"
	"net/http"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) HandleMiddleware(next http.HandlerFunc, mws ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	// We always want to log the request
	mws = append(mws, m.LogRequest)
	for _, mw := range mws {
		next = mw(next.ServeHTTP)
	}
	return next
}

func (m *Middleware) LogRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
