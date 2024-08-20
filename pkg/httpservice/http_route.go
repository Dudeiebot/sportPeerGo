package httpservice

import (
	"github.com/go-chi/chi/v5"

	"github.com/dudeiebot/sportPeerGo/pkg/user"
)

func AuthRoutes(r chi.Router, s *Server) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", user.AddHostSchemeMiddleware(CreateUser(s)))
		r.Post("/login", LoginUser(s))
		// r.Post("/logout", logoutUser(dbs))
		// r.Post("/verify-otp", verifyOtp(dbs))
		r.Get("/verify-email", VerifyEmail(s))
	})
}
