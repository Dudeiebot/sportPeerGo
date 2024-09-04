package smtps

import (
	"fmt"
	"net/smtp"
	"os"

	emailNew "github.com/jordan-wright/email"
)

type UserInfo struct {
	RecipientEmail    string
	VerificationToken string
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

func SendVerificationEmail(info *UserInfo, host, scheme string) error {
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

func SendOtpEmail(info *UserInfo, host, scheme string) error {
	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", config.FromEmail)
	e.To = []string{info.RecipientEmail}
	e.Subject = "Email Verification Link"
	e.Text = []byte(
		fmt.Sprintf(
			"Please Click the link to change your password: %s://%s/auth/pass?otptoken=%s?&email=%s",
			scheme,
			host,
			info.VerificationToken,
			info.RecipientEmail,
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
