package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	random "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// func createUser(dbService *dbs.Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user model.User
// 		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
//
// 		if err := user.ValidateUser(); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
//
// 		pwbytes := []byte(user.Password)
// 		hashedPass, err := bcrypt.GenerateFromPassword(pwbytes, bcrypt.DefaultCost)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		user.Password = string(hashedPass)
// 		user.Username = generateUsername(user.Email)
// 		user.VerificationToken, err = verificationToken()
// 		if err != nil {
// 			http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
// 			return
// 		}
//
// 		user.ID = queries.RegisterQuery(user, dbService)
// 		if user.ID == 0 {
// 			http.Error(w, "Failed to register user", http.StatusInternalServerError)
// 			return
// 		}
//
// 		info := &email.UserInfo{
// 			RecipientEmail:    user.Email,
// 			VerificationToken: user.VerificationToken,
// 			Req:               r,
// 		}
//
// 		go func() {
// 			err = email.SendVerificationEmail(context.Background(), info)
// 			if err != nil {
// 				log.Printf("Failed to send verification email: %v", err)
// 			}
// 		}()
//
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"message": "User registered successfully. Please check your email for the verification link.",
// 			"user":    user,
// 		})
// 	}
// }

func VerificationToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 1 {
		baseUsername := parts[0]
		r := random.New(random.NewSource(time.Now().Unix()))
		randNum := r.Intn(1000)
		return baseUsername + strconv.Itoa(randNum)
	}
	return ""
}

func EncryptAuth(pass string) (string, error) {
	hashedAuth, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedAuth), nil
}

func AddHostSchemeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		ctx := context.WithValue(r.Context(), "host", r.Host)
		ctx = context.WithValue(ctx, "scheme", scheme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func loginUser(dbService *dbs.Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userID := chi.URLParam(r, "userID")
// 		// Your logic to get a user by ID
// 		w.Write([]byte("Get User by ID: " + userID))
// 	}
// }
//
// func logoutUser(dbService *dbs.Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userID := chi.URLParam(r, "userID")
// 		// Your logic to update a user by ID
// 		w.Write([]byte("Update User by ID: " + userID))
// 	}
// }
//
// func verifyOtp(dbService *dbs.Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userID := chi.URLParam(r, "userID")
// 		// Your logic to delete a user by ID
// 		w.Write([]byte("Delete User by ID: " + userID))
// 	}
// }
