package user

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	random "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Claim struct {
	Subject   int64 `json:"sub"`
	IssuedAt  int64 `json:"iat"`
	ExpiredAt int64 `json:"exp"`
}

var secret []byte

func GenerateSecretToken(id int64) (string, error) {
	secret = []byte(os.Getenv("SECRET"))
	now := time.Now().Unix()

	claims := Claim{
		Subject:   id,
		IssuedAt:  now,
		ExpiredAt: now + 3600,
	}

	claimJson, err := json.Marshal(claims)
	if err != nil {
		return "", nil
	}
	payLoad := base64.URLEncoding.EncodeToString(claimJson)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(payLoad))
	sig := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := payLoad + "." + sig
	return token, nil
}

func VerificationToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 1 {
		baseUsername := parts[0]
		r := random.New(random.NewSource(time.Now().Unix()))
		randNum := r.Intn(1000)
		return baseUsername + strconv.Itoa(randNum)
	}
	return ""
}

func EncryptAuth(pass string) (string, error) {
	hashedAuth, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedAuth), nil
}

func AddHostSchemeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		ctx := context.WithValue(r.Context(), "host", r.Host)
		ctx = context.WithValue(ctx, "scheme", scheme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
