package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerateSecretToken(t *testing.T) {
	os.Setenv("SECRET", "test-secret")
	defer os.Unsetenv("SECRET")

	id := int64(123)
	token, err := GenerateSecretToken(id)
	if err != nil {
		t.Fatalf("GenerateSecretToken failed: %v", err)
	}
	if token == "" {
		t.Error("GenerateSecretToken returned empty token")
	}
}

func TestVerificationToken(t *testing.T) {
	token, err := VerificationToken()
	if err != nil {
		t.Fatalf("VerificationToken failed: %v", err)
	}
	if len(token) != 32 {
		t.Errorf(
			"VerificationToken returned token of unexpected length: got %d, want 32",
			len(token),
		)
	}
}

func TestGenerateUsername(t *testing.T) {
	email := "test@example.com"
	username := GenerateUsername(email)
	if !strings.HasPrefix(username, "test") {
		t.Errorf("GenerateUsername did not use email prefix: got %s, want prefix 'test'", username)
	}
	if len(username) <= 4 {
		t.Errorf("GenerateUsername did not append random number: got %s", username)
	}
}

func TestEncryptAndCompareAuth(t *testing.T) {
	auth := "password123"
	encrypted, err := EncryptAuth(auth)
	if err != nil {
		t.Fatalf("EncryptAuth failed: %v", err)
	}
	if encrypted == auth {
		t.Error("EncryptAuth did not change the password")
	}

	err = CompareAuth(encrypted, auth)
	if err != nil {
		t.Errorf("CompareAuth failed for correct password: %v", err)
	}

	err = CompareAuth(encrypted, "wrongpassword")
	if err == nil {
		t.Error("CompareAuth did not fail for incorrect password")
	}
}

func TestIsLoggedOut(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	if !IsLoggedOut(r) {
		t.Error("IsLoggedOut returned false for request without cookie")
	}

	r.AddCookie(&http.Cookie{Name: "token", Value: "test-token"})
	if IsLoggedOut(r) {
		t.Error("IsLoggedOut returned true for request with token cookie")
	}
}

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("SECRET", "test-secret")
	defer os.Unsetenv("SECRET")

	handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId")
		if userId == nil {
			t.Error("AuthMiddleware did not set userId in context")
		}
	})

	// Test with valid token
	token, _ := GenerateSecretToken(123)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("AuthMiddleware failed for valid token: got status %d", w.Code)
	}

	// Test with invalid token
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer invalidtoken")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("AuthMiddleware did not fail for invalid token: got status %d", w.Code)
	}
}

func TestValidateToken(t *testing.T) {
	os.Setenv("SECRET", "test-secret")
	defer os.Unsetenv("SECRET")

	// Generate a valid token
	claims := &Claim{
		Subject:   123,
		IssuedAt:  time.Now().Unix(),
		ExpiredAt: time.Now().Add(time.Hour).Unix(),
	}
	claimJson, _ := json.Marshal(claims)
	payload := base64.URLEncoding.EncodeToString(claimJson)
	h := hmac.New(sha256.New, []byte("test-secret"))
	h.Write([]byte(payload))
	sig := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := payload + "." + sig

	// Test valid token
	validatedClaims, err := ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken failed for valid token: %v", err)
	}
	if validatedClaims.Subject != claims.Subject {
		t.Errorf(
			"ValidateToken returned wrong subject: got %d, want %d",
			validatedClaims.Subject,
			claims.Subject,
		)
	}

	// Test invalid token
	_, err = ValidateToken("invalid.token")
	if err == nil {
		t.Error("ValidateToken did not fail for invalid token")
	}
}
