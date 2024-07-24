package server

import (
	"fmt"
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

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	dbService := dbs.New()
	ServerInstance := &Server{
		port: port,
		dbs:  dbService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	user.UserRoutes(r, dbService)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", ServerInstance.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
