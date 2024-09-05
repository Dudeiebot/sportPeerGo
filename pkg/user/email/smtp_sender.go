package smtps

import (
	"fmt"
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
	recipientEmail, subject string,
	textContent func(string, string) string,
	host, scheme string,
) error {
	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", config.FromEmail)
	e.To = []string{recipientEmail}
	e.Subject = subject

	e.Text = []byte(textContent(host, scheme))

	err := e.Send(
		fmt.Sprintf("%s:%s", config.SMTPServer, config.SMTPPort),
		smtp.PlainAuth("", config.PostmarkToken, config.PostmarkToken, config.SMTPServer),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func SendVerificationEmail(info *UserInfo, host, scheme string) error {
	return SendEmail(
		info.RecipientEmail,
		"Email Verification Link",
		func(host, scheme string) string {
			return fmt.Sprintf(
				"Please verify your email by clicking the link: %s://%s/auth/verify-email?token=%s",
				scheme, host, info.Token,
			)
		},
		host, scheme,
	)
}

func SendOtpEmail(info *UserInfo, host, scheme string) error {
	return SendEmail(
		info.RecipientEmail,
		"Password Changing Link",
		func(host, scheme string) string {
			return fmt.Sprintf(
				"Please Click the link to change your password: %s://%s/auth/updatepass?otptoken=%s&email=%s",
				scheme,
				host,
				info.Token,
				info.RecipientEmail,
			)
		},
		host, scheme,
	)
}
