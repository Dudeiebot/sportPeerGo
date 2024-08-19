package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validErrorMessages = map[string]map[string]string{
	"Email": {
		"required": "Email address is required",
		"email":    "Please enter a valid email address",
	},
	"Phone": {
		"required": "Phone number is required",
		"e164":     "Please enter a valid phone number in +12343535 format",
	},
	"Password": {
		"required": "Password is required",
		"gte":      "Password must be at least 6 characters long",
	},
}

func ValidatorErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessage := make([]string, 0, len(validationErrors))
		for _, e := range validationErrors {
			errorMessage = append(errorMessage, getValidErrorMessage(e))
		}
		return fmt.Errorf("validation failed:\n%s", strings.Join(errorMessage, "\n"))
	}
	return err
}

func getValidErrorMessage(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()

	if messages, ok := validErrorMessages[field]; ok {
		if message, ok := messages[tag]; ok {
			return message
		}
	}
	return "Write the right Input"
}
