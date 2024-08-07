package email

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	emailNew "github.com/jordan-wright/email"
)

type UserInfo struct {
	RecipientEmail    string
	VerificationToken string
	Req               *http.Request
}

func SendVerificationEmail(info *UserInfo) error {
	var (
		fromEmail     = os.Getenv("FROM")
		smtpServer    = os.Getenv("SMTP_SERVER")
		smtpPort      = os.Getenv("SMTP_PORT")
		postmarkToken = os.Getenv("POSTMARK_TOKEN")
	)

	host := info.Req.Host
	scheme := "http"
	if info.Req.TLS != nil {
		scheme = "https"
	}

	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", fromEmail)
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
		fmt.Sprintf("%s:%s", smtpServer, smtpPort),
		smtp.PlainAuth(
			"",
			postmarkToken,
			postmarkToken,
			smtpServer,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
