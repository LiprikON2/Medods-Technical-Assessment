package validator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	auth "github.com/medods-technical-assessment"
)

type ValidationService struct {
	validate *validator.Validate
}

func NewValidationService() *ValidationService {
	return &ValidationService{
		validate: InitValidator(),
	}
}

func InitValidator() *validator.Validate {
	validate := validator.New()

	// Register function to return json tags in field validation (e.g. `err.Field()`)
	// ref: https://github.com/go-playground/validator/issues/258#issuecomment-257281334
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom password validation
	validate.RegisterValidation("password", isValidPassword)

	return validate
}

// ref: https://stackoverflow.com/a/25840157
func isValidPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasLower,
		hasUpper,
		hasNumber,
		hasSpecial bool
	)
	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (s *ValidationService) ValidateUserInput(input any) []auth.ValidationError {
	var errors []auth.ValidationError

	err := s.validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string

			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				message = "invalid email format"
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			case "max":
				message = fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
			case "password":
				message = "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
			default:
				message = fmt.Sprintf("invalid value for %s", err.Field())
			}

			errors = append(errors, auth.ValidationError{
				Field:   err.Field(),
				Message: message,
			})
		}
	}

	return errors
}
