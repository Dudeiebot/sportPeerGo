package model

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/dudeiebot/sportPeerGo/pkg/errors"
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
	Email          string `json:"email"`
	Otp            string `json:"otp"`
	ExpirationTime time.Time
	NewPass        string `json:"password"`
}

func (u *User) ValidateUser() error {
	rules := map[string]string{
		"Email":    "required,email",
		"Phone":    "required,e164",
		"Password": "required,gte=6",
	}

	valid := validator.New()
	valid.RegisterStructValidationMapRules(rules, u)

	err := valid.Struct(u)
	if err != nil {
		return errors.ValidatorErrors(err)
	}
	return nil
}
