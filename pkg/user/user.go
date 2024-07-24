package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dudeiebot/sportPeerGo/pkg/dbs"
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
		// Your logic to create a user
		w.Write([]byte("Create User"))
	}
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
