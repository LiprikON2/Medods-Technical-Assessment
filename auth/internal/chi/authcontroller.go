package chi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/pkg/utils"
)

type AuthController struct {
	service auth.AuthService
}

func NewAuthController(service auth.AuthService) *AuthController {
	return &AuthController{
		service: service,
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

	userDto := user.ToDto()

	if err := encoder.Encode(userDto); err != nil {
		log.Print(err)
		InternalErrorHandler(w)
		return
	}
}

func (c *AuthController) Users(w http.ResponseWriter, r *http.Request) {

	users, err := c.service.Users()

	usersDto := utils.MapSlice(users, func(user *auth.User) auth.UserDto {
		return user.ToDto()
	})

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(usersDto); err != nil {
		log.Print(err)
		InternalErrorHandler(w)
		return
	}
}

var (
	BadRequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusNotFound)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpected Error Occured.", http.StatusInternalServerError)
	}
)

func writeError(w http.ResponseWriter, message string, statusCode int) {
	resp := auth.Error{
		Code:    statusCode,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)

	encoder.SetIndent("", "  ")
	encoder.Encode(resp)
}
