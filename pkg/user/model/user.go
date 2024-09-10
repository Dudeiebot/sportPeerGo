package model

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Password          string `json:"password"`
	VerificationToken string `json:"VerificationToken"`
	Bio               string `json:"bio"`
	IsVerified        bool   `json:"is_verified"`
}

type Credentials struct {
	Access   string `json:"access"`
	Password string `json:"password"`
}

type ForgetPass struct {
	Email          string    `json:"email"`
	Otp            string    `json:"otp"`
	ExpirationTime time.Time `json:"expiration_time"`
	NewPass        string    `json:"password"`
}

func (u *User) ValidateUser() error {
	var errors []string

	for field, value := range map[string]string{
		"Email":    u.Email,
		"Phone":    u.Phone,
		"Password": u.Password,
	} {
		switch field {
		case "Email":
			if value == "" {
				errors = append(errors, "Email is required")
			} else if !isValidEmail(value) {
				errors = append(errors, "Invalid email format")
			}
		case "Phone":
			if value == "" {
				errors = append(errors, "Phone is required")
			} else if !isValidPhone(value) {
				errors = append(errors, "Invalid phone number format")
			}
		case "Password":
			if value == "" {
				errors = append(errors, "Password is required")
			} else if len(value) < 6 {
				errors = append(errors, "Password must be at least 6 characters long")
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Validation errors: %s", strings.Join(errors, ", "))
	}
	return nil
}

func (c *Credentials) ValidateCred() error {
	var errors []string
	for field, value := range map[string]string{
		"Access":   c.Access,
		"Password": c.Password,
	} {
		switch field {
		case "Access":
			if value == "" {
				errors = append(errors, "Email or Phone is required")
			} else if !isValidEmail(c.Access) && !isValidPhone(c.Access) {
				errors = append(errors, "Access must be a valid email or phone number")
			}
		case "Password":
			if value == "" {
				errors = append(errors, "Password is required")
			} else if len(value) < 6 {
				errors = append(errors, "Password must be at least 6 characters long")
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("Validation errors: %s", strings.Join(errors, ", "))
	}
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPhone(phone string) bool {
	pattern := `^\+[1-9]\d{1,14}$`
	match, _ := regexp.MatchString(pattern, phone)
	return match
}
