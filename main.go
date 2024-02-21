package main

import (
	"log"
	"net/http"
)

func main() {
	conf, err := LoadConfig()
	if err != nil {
		log.Fatal("LoadConfig: ", err)
	}

	h := NewHandlers(conf)
	r := NewRouter(h)

	log.Printf("Server started on http://localhost:8080\n")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
