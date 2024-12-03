package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	auth "github.com/medods-technical-assessment"
)

type Handler struct {
	service auth.AuthService
}

func NewHandler(service auth.AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListenAndServe(r *chi.Mux) {

	log.Print("Starting HTTP server at: http://localhost:8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle request
	log.Print("ResponseWriter", w, "request", r)

}
