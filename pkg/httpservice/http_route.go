package httpservice

import (
	"github.com/go-chi/chi/v5"

	"github.com/dudeiebot/sportPeerGo/pkg/user"
)

func AuthRoutes(r chi.Router, s *Server) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", user.AddHostSchemeMiddleware(CreateUser(s)))
		r.Post("/login", LoginUser(s))
		r.Post("/logout", LogoutUser(s))
		// r.Post("/verify-otp", verifyOtp(dbs))
		r.Get("/verify-email", VerifyEmail(s))
	})
}

func UserRoute(r chi.Router, s *Server) {
	r.Route("/users", func(r chi.Router) {
		r.Put("/username/{id}", user.AuthMiddleware(UpdateUsername(s)))
		r.Put("/email/{id}", user.AuthMiddleware(UpdateEmail(s)))
	})
}
