package model

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/dudeiebot/sportPeerGo/pkg/user/errors"
)

type User struct {
	ID                int       `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Password          string    `json:"password"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updated_at"`
	VerificationToken string    `json:"verification_token"`
	IsVerified        bool      `json:"is_verified"`
	VerificationID    string    `json:"verification_id"`
	Bio               string    `json:"bio"`
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
