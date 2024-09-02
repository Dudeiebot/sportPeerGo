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
	"golang.org/x/crypto/bcrypt"

	"github.com/dudeiebot/sportPeerGo/pkg/adapter/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/adapter/queries"
	"github.com/dudeiebot/sportPeerGo/pkg/user"
	"github.com/dudeiebot/sportPeerGo/pkg/user/email"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
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

func CreateUser(s *Server) http.HandlerFunc {
	return NewHandler(
		func(ctx context.Context, u model.User) (map[string]interface{}, error) {
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

func VerifyEmail(s *Server) http.HandlerFunc {
	return NewHandler(func(ctx context.Context, r *http.Request) (*Response, error) {
		token := r.URL.Query().Get("token")

		res, err := queries.VerifyEmailQueries(ctx, s.DBS, token)
		if err != nil {
			log.Printf("Error executing db query: %v\n", err)
			return nil, fmt.Errorf("internal server error")
		}

		rowAffected, err := res.RowsAffected()

		if rowAffected == 0 {
			return &Response{Message: "Invalid or expired token"}, nil
		}

		return &Response{Message: "Email Verifed Successfully"}, nil
	})
}

func LoginUser(s *Server) http.HandlerFunc {
	return NewHandler(
		func(ctx context.Context, c model.Credentials) (*LoginResponse, error) {
			u, err := queries.GetHashedAuth(ctx, c, s.DBS)
			if err != nil {
				return nil, fmt.Errorf(
					"Invalid Credentials, Please provide the correct email or phone number",
				)
			}

			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Password)); err != nil {
				return nil, fmt.Errorf("Invalid Credentials, Please provide the correct password")
			}

			if !u.IsVerified {
				return nil, fmt.Errorf("Please Verify your account before logging in")
			}

			accessToken, err := user.GenerateSecretToken(int64(u.ID))
			if err != nil {
				return nil, err
			}

			u.Password = ""

			return &LoginResponse{
				Message: "User Logged In Successfully",
				Token:   accessToken,
			}, nil
		},
	)
}

func LogoutUser(s *Server) http.HandlerFunc {
	return NewHandler(func(ctx context.Context, r *http.Request) (*LogoutResponse, error) {
		return &LogoutResponse{
			Message: "User Logged Out",
		}, nil
	})
}

func UpdateUsername(s *Server) http.HandlerFunc {
	return NewUpdateHandler(s, queries.UsernameQueries, "Username successfully changed")
}

func UpdateEmail(s *Server) http.HandlerFunc {
	return NewUpdateHandler(s, queries.EmailQueries, "Email changed successfully")
}
