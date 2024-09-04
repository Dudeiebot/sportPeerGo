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

func EncryptAuth(auth string) (string, error) {
	hashedAuth, err := bcrypt.GenerateFromPassword([]byte(auth), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedAuth), nil
}

func CompareAuth(hashedAuth, unhashedAuth string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedAuth), []byte(unhashedAuth))
	if err != nil {
		return err
	}
	return nil
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

func IsLoggedOut(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		return true
	}

	return false
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if token == authHeader {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", claims.Subject)
		next(w, r.WithContext(ctx))
	}
}

func ValidateToken(token string) (*Claim, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload := parts[0]
	sig := parts[1]

	secret := []byte(os.Getenv("SECRET"))
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(payload))
	expectedSig := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if sig != expectedSig {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimJson, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding")
	}

	var claims Claim
	if err := json.Unmarshal(claimJson, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims.ExpiredAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func GenerateOTP() (string, error) {
	randomBytes := make([]byte, 3)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	dotp := int64(randomBytes[0])<<16 | int64(randomBytes[1])<<8 | int64(randomBytes[2])
	dotp = dotp % 1000000
	return strconv.FormatInt(dotp, 10), nil
}
