package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dudeiebot/sportPeerGo/pkg/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user"
)

type Server struct {
	port int
	dbs  *dbs.Service
}

func NewServer(ctx context.Context) (*http.Server, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	dbService := dbs.New(ctx)
	serverInstance := &Server{
		port: port,
		dbs:  dbService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	user.UserRoutes(r, serverInstance.dbs)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverInstance.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server initialized, listening on port %d", serverInstance.port)
	return server, nil
}
