package model

import (
	"github.com/go-playground/validator/v10"

	"github.com/dudeiebot/sportPeerGo/pkg/user/errors"
)

type User struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Password          string `json:"password"`
	VerificationToken string `json:"VerificationToken"`
	Bio               string `json:"bio"`
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
