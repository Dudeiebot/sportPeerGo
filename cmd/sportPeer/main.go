package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dudeiebot/sportPeerGo/pkg/httpservice"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := httpservice.NewServer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
