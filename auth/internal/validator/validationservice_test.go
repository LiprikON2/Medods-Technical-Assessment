package validator

import (
	"testing"

	auth "github.com/medods-technical-assessment"
)

func TestValidationService(t *testing.T) {

	var tests = []struct {
		name      string
		input     *auth.UpdateUserDto
		wantValid bool
	}{
		{"Valid email and valid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "Hello1234!",
		}, true},
		{"Valid email and valid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "Hello1234!",
		}, true},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "Hello1234",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "Hel14!",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "hello1234!",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "hello1234",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "Hello1234",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "HELLO1234",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "HELLOhello",
		}, false},
		{"Valid email and invalid password", &auth.UpdateUserDto{
			Email:    "test@example.com",
			Password: "HELLOhello!",
		}, false},
		{"Invalid email and valid password", &auth.UpdateUserDto{
			Email:    "test@example.co1",
			Password: "Hello1234!",
		}, false},
		{"Invalid email and valid password", &auth.UpdateUserDto{
			Email:    "test@.com",
			Password: "Hello1234!",
		}, false},
		{"Invalid email and valid password", &auth.UpdateUserDto{
			Email:    "test@example@s.com",
			Password: "Hello1234!",
		}, false},
		{"Invalid email and valid password", &auth.UpdateUserDto{
			Email:    "@example@s.com",
			Password: "Hello1234!",
		}, false},
		{"Invalid email and valid password", &auth.UpdateUserDto{
			Email:    "@example.com",
			Password: "Hello1234!",
		}, false},
		{"Invalid email and invalid password", &auth.UpdateUserDto{
			Email:    "@example.com",
			Password: "Hello124",
		}, false},
		{"Missing email and valid password", &auth.UpdateUserDto{
			Password: "Hello124",
		}, false},
		{"Valid email and missing password", &auth.UpdateUserDto{
			Email: "@example.com",
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := NewValidationService()
			ans := vs.ValidateUserInput(tt.input)
			isValid := ans == nil

			if isValid != tt.wantValid {
				t.Errorf("got valid %v, want valid %v", isValid, tt.wantValid)
			}
		})
	}
}
