package chi

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/medods-technical-assessment"
)

var (
	ValidationErrorHandler = func(w http.ResponseWriter, errors []auth.ValidationError) {
		message := map[string]any{
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
	ForbiddenErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusForbidden)
	}
	InternalErrorHandler = func(w http.ResponseWriter, err error) {
		log.Print(err)
		writeError(w, "An Unexpected Error Occured.", http.StatusInternalServerError)
	}
)

func writeError(w http.ResponseWriter, message any, statusCode int) {
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
