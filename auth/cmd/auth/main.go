package main

import (
	"log"
	"time"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	"github.com/medods-technical-assessment/internal/http"
	"github.com/medods-technical-assessment/internal/postgres"
	tables "github.com/medods-technical-assessment/internal/postgres/tables"
)

func main() {

	// Connect to database.
	db, err := postgres.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if credentials are valid
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	if err := tables.CreateUsersTable(db); err != nil {
		log.Fatal(err)
	}

	// Create services.
	as := &postgres.AuthService{DB: db}

	// Attach to HTTP handler.
	var h = http.NewHandler(as)

	var c = postgres.NewAuthController(as)

	// start http server...
	var r *chi.Mux = chi.NewRouter()

	r.Use(chimiddle.StripSlashes)

	// A good base middleware stack
	r.Use(chimiddle.RequestID)
	r.Use(chimiddle.RealIP)
	r.Use(chimiddle.Logger)
	r.Use(chimiddle.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(chimiddle.Timeout(60 * time.Second))

	r.Route("/auth", func(r chi.Router) {
		// router.Use(middleware.Authorization)

		r.Route("/", func(r chi.Router) {
			// r.Use(c.UserCtx)
			r.Get("/{UserID}", c.GetUser)
			// r.Put("/", updateArticle)                                       // PUT /articles/123
			// r.Delete("/", deleteArticle)                                    // DELETE /articles/123
		})
	})

	// handlers.Handler(r)
	h.ListenAndServe(r)
}
