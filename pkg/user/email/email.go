package email

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	emailNew "github.com/jordan-wright/email"
)

func SendVerificationEmail(recipientEmail, verificationToken string, req *http.Request,
) error {
	var (
		fromEmail     = os.Getenv("FROM")
		smtpServer    = os.Getenv("SMTP_SERVER")
		smtpPort      = os.Getenv("SMTP_PORT")
		postmarkToken = os.Getenv("POSTMARK_TOKEN")
	)

	host := req.Host
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}

	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", fromEmail)
	e.To = []string{recipientEmail}
	e.Subject = "Email Verification Link"
	e.Text = []byte(
		fmt.Sprintf(
			"Please verify your email by clicking the link: %s://%s/auth/verify-email?token=%s",
			scheme,
			host,
			verificationToken,
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
