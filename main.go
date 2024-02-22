package main

import (
	"log"
	"net/http"
)

func main() {
	c, err := LoadConfig()
	if err != nil {
		log.Fatal("LoadConfig: ", err)
	}

	s := NewService(c)
	h := NewHandlers(c, s)
	mw := NewMiddleware()
	r := NewRouter(mw, h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Printf("Server listening on %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
