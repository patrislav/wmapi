package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"

	"github.com/patrislav/wmapi/rest"
	"github.com/patrislav/wmapi/x11"
)

func main() {
	xclient, err := x11.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize an X client: %v", err)
	}

	api := rest.NewServer(xclient)
	server := chi.NewRouter()
	server.Mount("/v1", api)

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Fatal(http.ListenAndServe(addr, server))
}
