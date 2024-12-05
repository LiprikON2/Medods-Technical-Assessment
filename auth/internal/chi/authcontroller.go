package chi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/internal/common"
)

type AuthController struct {
	service           auth.AuthService
	validationService auth.ValidationService
}

func NewAuthController(service auth.AuthService, validationService auth.ValidationService) *AuthController {
	return &AuthController{
		service:           service,
		validationService: validationService,
	}
}

func (c *AuthController) User(w http.ResponseWriter, r *http.Request) {

	userUUID, ok := r.Context().Value("UserUUID").(uuid.UUID)
	if !ok {
		InternalErrorHandler(w, fmt.Errorf("failed to get UUID from context"))
		return
	}
	user, err := c.service.User(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(user); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) Users(w http.ResponseWriter, r *http.Request) {

	users, err := c.service.Users()

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(users); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) CreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userInput auth.CreateUserDto
	if err := decoder.Decode(&userInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}

	// Validate input
	if errors := c.validationService.ValidateUserInput(userInput); len(errors) > 0 {
		ValidationErrorHandler(w, errors)
		return
	}

	user := &auth.User{
		UUID:  uuid.New(),
		Email: userInput.Email,
		// TODO hash
		Password: userInput.Password,
	}

	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateEmail) {
			ConflictErrorHandler(w, err)
			return
		}
		InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(createdUser); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userUUID, ok := r.Context().Value("UserUUID").(uuid.UUID)
	if !ok {
		InternalErrorHandler(w, fmt.Errorf("failed to get UUID from context"))
		return
	}
	user, err := c.service.User(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userInput auth.UpdateUserDto
	if err := decoder.Decode(&userInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}

	// Validate input
	if errors := c.validationService.ValidateUserInput(userInput); len(errors) > 0 {
		ValidationErrorHandler(w, errors)
		return
	}

	if userInput.Email != "" {
		user.Email = userInput.Email
	}
	// TODO hash
	if userInput.Password != "" {
		user.Password = userInput.Password
	}

	updatedUser, err := c.service.UpdateUser(user)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateEmail) {
			ConflictErrorHandler(w, err)
			return
		}
		InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(updatedUser); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userUUID, ok := r.Context().Value("UserUUID").(uuid.UUID)
	if !ok {
		InternalErrorHandler(w, fmt.Errorf("failed to get UUID from context"))
		return
	}
	err := c.service.DeleteUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
