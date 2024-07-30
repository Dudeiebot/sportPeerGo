package model

import "time"

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
