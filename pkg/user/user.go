package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	random "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/dudeiebot/sportPeerGo/pkg/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user/email"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
	"github.com/dudeiebot/sportPeerGo/pkg/user/queries"
)

func UserRoutes(r chi.Router, dbs *dbs.Service) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", createUser(dbs))
		r.Post("/login", loginUser(dbs))
		r.Post("/logout", logoutUser(dbs))
		r.Post("/verify-otp", verifyOtp(dbs))
		r.Post("/verify-email", verifyEmail(dbs))
	})
}

func createUser(dbService *dbs.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := user.ValidateUser(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pwbytes := []byte(user.Password)
		hashedPass, err := bcrypt.GenerateFromPassword(pwbytes, bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedPass)
		user.Username = generateUsername(user.Email)
		user.VerificationToken, err = verificationToken()
		if err != nil {
			http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
			return
		}

		user.ID = queries.RegisterQuery(user, dbService)
		if user.ID == 0 {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		go func() {
			err = email.SendVerificationEmail(user.Email, user.VerificationToken, r)
			if err != nil {
				log.Printf("Failed to send verification email: %v", err)
				// Decide if you want to return an error to the client or just log it
				// http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
				// return
			}
		}()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User registered successfully. Please check your email for the verification link.",
			"user":    user,
		})
	}
}

func verificationToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func generateUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 1 {
		baseUsername := parts[0]
		r := random.New(random.NewSource(time.Now().Unix()))
		randNum := r.Intn(1000)
		return baseUsername + strconv.Itoa(randNum)
	}
	return ""
}

func loginUser(dbService *dbs.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		// Your logic to get a user by ID
		w.Write([]byte("Get User by ID: " + userID))
	}
}

func logoutUser(dbService *dbs.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		// Your logic to update a user by ID
		w.Write([]byte("Update User by ID: " + userID))
	}
}

func verifyOtp(dbService *dbs.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		// Your logic to delete a user by ID
		w.Write([]byte("Delete User by ID: " + userID))
	}
}

func verifyEmail(dbService *dbs.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		// Your logic to delete a user by ID
		w.Write([]byte("Delete User by ID: " + userID))
	}
}
