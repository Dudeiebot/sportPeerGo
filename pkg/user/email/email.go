package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	emailNew "github.com/jordan-wright/email"

	"github.com/dudeiebot/sportPeerGo/pkg/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user/queries"
)

type UserInfo struct {
	RecipientEmail    string
	VerificationToken string
	Req               *http.Request
}

type Response struct {
	Message string `json:"message"`
}

type Config struct {
	FromEmail     string
	SMTPServer    string
	SMTPPort      string
	PostmarkToken string
}

var config = &Config{
	FromEmail:     os.Getenv("FROM"),
	SMTPServer:    os.Getenv("SMTP_SERVER"),
	SMTPPort:      os.Getenv("SMTP_PORT"),
	PostmarkToken: os.Getenv("POSTMARK_TOKEN"),
}

func SendVerificationEmail(ctx context.Context, info *UserInfo) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	host := info.Req.Host
	scheme := "http"
	if info.Req.TLS != nil {
		scheme = "https"
	}

	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", config.FromEmail)
	e.To = []string{info.RecipientEmail}
	e.Subject = "Email Verification Link"
	e.Text = []byte(
		fmt.Sprintf(
			"Please verify your email by clicking the link: %s://%s/auth/verify-email?token=%s",
			scheme,
			host,
			info.VerificationToken,
		),
	)

	err := e.Send(
		fmt.Sprintf("%s:%s", config.SMTPServer, config.SMTPPort),
		smtp.PlainAuth(
			"",
			config.PostmarkToken,
			config.PostmarkToken,
			config.SMTPServer,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func VerifyEmail(dbService *dbs.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.URL.Query().Get("token")

		res, err := queries.VerifyEmailQueries(ctx, dbService, token)
		if err != nil {
			log.Printf("Error executing db query: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		rowAffected, err := res.RowsAffected()

		if rowAffected == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Invalid or expired token"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Message: "Email Verifed Successfully"})
	}
}
