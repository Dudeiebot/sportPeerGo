package smtps

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	emailNew "github.com/jordan-wright/email"
)

type UserInfo struct {
	RecipientEmail string
	Token          string
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

func SendEmail(
	recipientEmail, subject, content string,
	r *http.Request,
) error {
	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", config.FromEmail)
	e.To = []string{recipientEmail}
	e.Subject = subject

	e.Text = []byte(content)
	err := e.Send(
		fmt.Sprintf("%s:%s", config.SMTPServer, config.SMTPPort),
		smtp.PlainAuth("", config.PostmarkToken, config.PostmarkToken, config.SMTPServer),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func SendVerificationEmail(info *UserInfo, r *http.Request) error {
	content := fmt.Sprintf(
		"Please verify your email by clicking the link: %s://%s/auth/verify-email?token=%s",
		getScheme(r), r.Host, info.Token,
	)
	return SendEmail(info.RecipientEmail, "Email Verification Link", content, r)
}

func SendOtpEmail(info *UserInfo, r *http.Request) error {
	content := fmt.Sprintf(
		"Please Click the link to change your password: %s://%s/auth/updatepass?otptoken=%s&email=%s",
		getScheme(r),
		r.Host,
		info.Token,
		info.RecipientEmail,
	)
	return SendEmail(info.RecipientEmail, "Password Changing Link", content, r)
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
