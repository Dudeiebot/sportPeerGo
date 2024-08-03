package main

import (
	"log"
	"net/http"

	"github.com/dudeiebot/sportPeerGo/pkg/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
