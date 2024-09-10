package model

import (
	"testing"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name string
		u    User
		want bool
	}{
		{
			name: "Valid user", u: User{Email: "test@example.com", Phone: "+1234567890", Password: "password123"}, want: false,
		},
		{
			name: "Invalid email", u: User{Email: "invalid-email", Phone: "+1234567890", Password: "password123"}, want: true,
		},
		{
			name: "Invalid phone", u: User{Email: "test@example.com", Phone: "invalid-phone", Password: "password123"}, want: true,
		},
		{
			name: "Short password", u: User{Email: "test@example.com", Phone: "+1234567890", Password: "short"}, want: true,
		},
		{
			name: "Empty fields", u: User{}, want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.u.ValidateUser()
			if (err != nil) != tt.want {
				t.Errorf("ValidateUser() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestValidateCred(t *testing.T) {
	tests := []struct {
		name string
		cred Credentials
		want bool
	}{
		{"Valid Email", Credentials{Access: "test@example.com", Password: "password123"}, false},
		{"Valid Phone", Credentials{Access: "+1234567890", Password: "password123"}, false},
		{"Empty Access", Credentials{Access: "", Password: "password123"}, true},
		{"Invalid Access", Credentials{Access: "invalid", Password: "password123"}, true},
		{"Short Password", Credentials{Access: "test@example.com", Password: "short"}, true},
		{"Empty Password", Credentials{Access: "test@example.com", Password: ""}, true},
		{"Empty Both", Credentials{Access: "", Password: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cred.ValidateCred()
			if (err != nil) != tt.want {
				t.Errorf("ValidateCred() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Valid email", "test@example.com", true},
		{"Invalid email - no @", "testexample.com", false},
		{"Invalid email - no domain", "test@", false},
		{"Invalid email - spaces", "test @example.com", false},
		{"Empty email", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"Valid phone", "+1234567890", true},
		{"Invalid phone - no +", "1234567890", false},
		{"Invalid phone - letters", "+1234abcd", false},
		{"Invalid phone - too long", "+123456789012345678", false},
		{"Empty phone", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPhone(tt.phone)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
