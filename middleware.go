package main

import (
	"github.com/jwtly10/CodeEcho-Server/logger"
	"net/http"
	"time"
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
		start := time.Now()
		log := logger.Get()
		defer func() {
			log.
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("incoming request")

		}()

		next.ServeHTTP(w, r)
	}
}
