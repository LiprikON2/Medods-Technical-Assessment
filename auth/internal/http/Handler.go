package http

import (
	"log"
	"net/http"

	auth "github.com/medods-technical-assessment"
)

type Handler struct {
	AuthService auth.AuthService
}

func (h *Handler) ListenAndServe() {

	log.Print("localhost:8080")
	http.ListenAndServe("localhost:8080", h)
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle request
	log.Print("ResponseWriter", w, "request", r)

}
