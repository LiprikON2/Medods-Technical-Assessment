package chi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

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

	userID := chi.URLParam(r, "UserID")
	log.Printf("GetUser ID: %s", userID)

	id, err := strconv.Atoi(userID)
	if err != nil {
		BadRequestErrorHandler(w, fmt.Errorf("invalid user ID format (%v): %w", userID, err))
		return
	}
	user, err := c.service.User(id)

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
	if errors := c.validationService.ValidateCreateUser(userInput); len(errors) > 0 {
		ValidationErrorHandler(w, errors)
		return
	}

	user := &auth.User{
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

var (
	ValidationErrorHandler = func(w http.ResponseWriter, errors []auth.ValidationError) {
		message := map[string]interface{}{
			"errors": errors,
		}

		writeError(w, message, http.StatusUnprocessableEntity)
	}
	ConflictErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusConflict)
	}
	BadRequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusNotFound)
	}
	InternalErrorHandler = func(w http.ResponseWriter, err error) {
		log.Print(err)
		writeError(w, "An Unexpected Error Occured.", http.StatusInternalServerError)
	}
)

func writeError(w http.ResponseWriter, message interface{}, statusCode int) {
	resp := auth.RequestError{
		Code:    statusCode,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)

	encoder.SetIndent("", "  ")
	encoder.Encode(resp)
}
