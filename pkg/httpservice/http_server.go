package httpservice

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

	"github.com/dudeiebot/sportPeerGo/pkg/adapter/dbs"
)

type Server struct {
	port int
	DBS  *dbs.Service
}

type Response struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"accessToken"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

func NewServer(ctx context.Context) (*http.Server, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	dbService := dbs.New(ctx)
	serverInstance := &Server{
		port: port,
		DBS:  dbService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	AuthRoutes(r, serverInstance)
	UserRoute(r, serverInstance)

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
