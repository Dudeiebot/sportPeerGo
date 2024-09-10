package httpservice

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

	"github.com/dudeiebot/sportPeerGo/pkg/adapter/dbs"
	query "github.com/dudeiebot/sportPeerGo/pkg/adapter/queries"
	"github.com/dudeiebot/sportPeerGo/pkg/user"
	smtps "github.com/dudeiebot/sportPeerGo/pkg/user/email"
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

func SendOtp(s *Server) http.HandlerFunc {
	return NewHandler(
		func(ctx context.Context, req *http.Request) (*Response, error) {
			email := chi.URLParam(req, "email")
			var f model.ForgetPass

			if req.ContentLength > 0 {
				if err := json.NewDecoder(req.Body).Decode(&f); err != nil {
					return nil, nil
				}
			}

			f.Email = email
			otp, err := user.GenerateOTP()
			if err != nil {
				return nil, fmt.Errorf("failed to generate OTP: %w", err)
			}
			hashedOtp, err := user.EncryptAuth(otp)
			if err != nil {
				return nil, err
			}
			f.Otp = hashedOtp
			f.ExpirationTime = time.Now().Add(5 * time.Minute)

			err = query.StoreOtpQuery(ctx, s.DBS, f)
			if err != nil {
				return nil, fmt.Errorf("failed to add user OTP: %w", err)
			}

			info := &smtps.UserInfo{
				RecipientEmail: f.Email,
				Token:          otp,
			}
			go func() {
				if err := smtps.SendOtpEmail(info, req); err != nil {
					log.Printf("Failed to send OTP Email: %v", err)
				}
			}()

			return &Response{Message: "Forget Password Link Sent Successfully"}, nil
		},
	)
}

func CreateUser(s *Server) http.HandlerFunc {
	return NewHandler(
		func(ctx context.Context, r *http.Request) (map[string]interface{}, error) {
			var u model.User
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				return nil, err
			}
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

			u.ID = query.RegisterQuery(u, s.DBS)
			if u.ID == 0 {
				return nil, fmt.Errorf("failed to register user and retrieve user ID")
			}

			info := &smtps.UserInfo{
				RecipientEmail: u.Email,
				Token:          u.VerificationToken,
			}

			go func() {
				if err := smtps.SendVerificationEmail(info, r); err != nil {
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

		res, err := query.VerifyEmailQuery(ctx, s.DBS, token)
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

func VerifyOtpAndUpdatePass(s *Server) http.HandlerFunc {
	return NewHandler(func(ctx context.Context, req *http.Request) (*Response, error) {
		otp := req.URL.Query().Get("otptoken")
		email := req.URL.Query().Get("email")
		forgetPass, err := query.GetOtpQuery(ctx, s.DBS, email)
		if err != nil {
			return nil, fmt.Errorf("error retrieving OTP info: %w", err)
		}

		if err := user.CompareAuth(forgetPass.Otp, otp); err != nil {
			return nil, fmt.Errorf("invalid OTP")
		}

		if time.Now().After(forgetPass.ExpirationTime) {
			return nil, fmt.Errorf("OTP has expired")
		}

		var f model.ForgetPass
		if err := json.NewDecoder(req.Body).Decode(&f); err != nil {
			return nil, err
		}

		hashedPass, err := user.EncryptAuth(f.NewPass)
		if err != nil {
			return nil, err
		}

		f.Email = email
		f.NewPass = hashedPass

		if err := query.UpdatePasswordQuery(ctx, s.DBS, f); err != nil {
			return nil, fmt.Errorf("error updating password: %w", err)
		}

		if err := query.ClearOtpQuery(ctx, s.DBS, email); err != nil {
			return nil, fmt.Errorf("error clearing OTP: %w", err)
		}

		return &Response{Message: "Password updated successfully"}, nil
	})
}

func LoginUser(s *Server) http.HandlerFunc {
	return NewHandler(
		func(ctx context.Context, c model.Credentials) (*LoginResponse, error) {
			u, err := query.GetHashedAuth(ctx, c, s.DBS)
			if err != nil {
				return nil, fmt.Errorf(
					"Invalid Credentials, Please provide the correct email or phone number",
				)
			}
			if err := user.CompareAuth(u.Password, c.Password); err != nil {
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
	return NewUpdateHandler(s, query.UsernameQuery, "Username successfully changed")
}

func UpdateEmail(s *Server) http.HandlerFunc {
	return NewUpdateHandler(s, query.EmailQuery, "Email changed successfully")
}
