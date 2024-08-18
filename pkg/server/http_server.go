package server

import (
	"context"
	"encoding/json"
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
	"github.com/dudeiebot/sportPeerGo/pkg/user/email"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
	"github.com/dudeiebot/sportPeerGo/pkg/user/queries"
)

type Server struct {
	port int
	DBS  *dbs.Service
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

func AuthRoutes(r chi.Router, s *Server) {
	r.Route("/auth", func(r chi.Router) {
		createUser := CreateUser(s)
		r.Post("/register", user.AddHostSchemeMiddleware(createUser))
		// r.Post("/login", loginUser(dbs))
		// r.Post("/logout", logoutUser(dbs))
		// r.Post("/verify-otp", verifyOtp(dbs))
		r.Get("/verify-email", email.VerifyEmail(s.DBS))
	})
}

func NewHandler[IN, OUT any](
	s *Server,
	targetFunc func(context.Context, *Server, IN) (OUT, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in IN
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out, err := targetFunc(r.Context(), s, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}
}

func CreateUser(s *Server) http.HandlerFunc {
	return NewHandler(
		s,
		func(ctx context.Context, s *Server, u model.User) (map[string]interface{}, error) {
			if err := u.ValidateUser(); err != nil {
				return nil, err
			}

			hashedPassword, err := user.EncryptAuth(u.Password)
			if err != nil {
				return nil, err
			}
			u.Password = hashedPassword

			u.Username = user.GenerateUsername(u.Email)

			u.VerificationToken, err = user.VerificationToken()
			if err != nil {
				return nil, err
			}

			u.ID = queries.RegisterQuery(u, s.DBS)
			if u.ID == 0 {
				return nil, fmt.Errorf("failed to register user and retrieve user ID")
			}

			info := &email.UserInfo{
				RecipientEmail:    u.Email,
				VerificationToken: u.VerificationToken,
			}

			host, _ := ctx.Value("host").(string)
			scheme, _ := ctx.Value("scheme").(string)

			go func() {
				if err := email.SendVerificationEmail(context.Background(), info, host, scheme); err != nil {
					log.Printf("Failed to send verification email: %v", err)
				}
			}()

			// Create response without password
			response := map[string]interface{}{
				"message": "User registered successfully. Please check your email for the verification link.",
				"user": map[string]interface{}{
					"id":       u.ID,
					"username": u.Username,
					"email":    u.Email,
				},
			}

			return response, nil
		},
	)
}
