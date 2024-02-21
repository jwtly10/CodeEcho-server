package main

import (
	"log"
	"net/http"
)

func main() {
	h := NewHandlers()
	r := NewRouter(h)

	log.Printf("Server started on http://localhost:8080\n")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
