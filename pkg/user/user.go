package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	random "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

func AddTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
