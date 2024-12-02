package http

import (
	"net/http"

	auth "github.com/medods-technical-assessment"
)

type Handler struct {
	AuthService auth.AuthService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle request
}
