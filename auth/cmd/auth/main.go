package main

import (
	"log"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	"github.com/medods-technical-assessment/internal/handlers"
	"github.com/medods-technical-assessment/internal/http"
	"github.com/medods-technical-assessment/internal/postgres"
)

func main() {

	// Connect to database.
	db, err := postgres.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create services.
	as := &postgres.AuthService{DB: db}

	// Attach to HTTP handler.
	var h http.Handler
	h.AuthService = as

	// start http server...
	var r *chi.Mux = chi.NewRouter()

	r.Use(chimiddle.StripSlashes)
	r.Use(chimiddle.Logger)
	r.Route("/auth", func(router chi.Router) {
		// router.Use(middleware.Authorization)

		router.Get("/", handlers.GetUser)
	})

	// handlers.Handler(r)
	h.ListenAndServe(r)
}
