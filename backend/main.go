package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	api := NewAPI(NewMemoryStore())

	log.Printf("backend listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, api.Routes()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
