package main

import (
	"fmt"
	"github.com/jwtly10/CodeEcho-Server/logger"
	"net/http"
)

func main() {
	c, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config")
		return
	}
	log := logger.Get()

	s := NewService(c)
	h := NewHandlers(c, s)
	mw := NewMiddleware()
	r := NewRouter(mw, h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Info().Str("address", server.Addr).Msg("Server listening")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
