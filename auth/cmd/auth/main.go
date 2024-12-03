package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/medods-technical-assessment/internal/chi"
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
	as := postgres.NewAuthService(db)
	r := chi.NewChiRouter()

	ac := chi.NewAuthController(as)

	r.Use(middleware.StripSlashes)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/auth", func(r chi.Router) {
		// router.Use(middleware.Authorization)

		r.Route("/", func(r chi.Router) {
			// r.Use(c.UserCtx)
			r.Get("/{UserID}", ac.User)
			// r.Put("/", updateArticle)                                       // PUT /articles/123
			// r.Delete("/", deleteArticle)                                    // DELETE /articles/123
		})
	})

	// Start http server...
	log.Print("Starting HTTP server at: http://localhost:8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
