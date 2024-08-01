package email

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	emailNew "github.com/jordan-wright/email"
)

func sendEmail(email, verificationToken string, r *http.Request, w http.ResponseWriter) {
	e := emailNew.NewEmail()
	e.From = fmt.Sprintf("<%s>", os.Getenv("FROM"))
	e.To = []string{email}
	e.Subject = "Email Verification Link"
	e.Text = []byte(
		fmt.Sprintf(
			"Please verify your email by clicking the link: %s://%s/auth/verify-email?token=%s",
			// verify route later
			r.URL.Scheme,
			r.Host,
			verificationToken,
		),
	)

	err := e.Send(
		"smtp.postmarkapp.com:587",
		smtp.PlainAuth(
			"",
			os.Getenv("POSTMARK_TOKEN"),
			os.Getenv("POSTMARK_TOKEN"),
			"smtp.postmarkapp.com",
		),
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(
		w,
		"User registered successfully. Please check your email for the verification link.",
	)
}
