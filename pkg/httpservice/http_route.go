package httpservice

import (
	"github.com/go-chi/chi/v5"

	"github.com/dudeiebot/sportPeerGo/pkg/user"
)

func AuthRoutes(r chi.Router, s *Server) {
	r.Route("/auth", func(r chi.Router) {
		createUser := CreateUser(s)
		r.Post("/register", user.AddHostSchemeMiddleware(createUser))
		// r.Post("/login", loginUser(dbs))
		// r.Post("/logout", logoutUser(dbs))
		// r.Post("/verify-otp", verifyOtp(dbs))
		verifyEmail := VerifyEmail(s)
		r.Get("/verify-email", verifyEmail)
	})
}
