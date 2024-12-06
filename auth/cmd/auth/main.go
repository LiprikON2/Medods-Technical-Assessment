package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/medods-technical-assessment/internal/bcrypt"
	"github.com/medods-technical-assessment/internal/chi"
	chimiddleware "github.com/medods-technical-assessment/internal/chi/middleware"
	"github.com/medods-technical-assessment/internal/jwt"
	"github.com/medods-technical-assessment/internal/postgres"
	"github.com/medods-technical-assessment/internal/uuid"
	"github.com/medods-technical-assessment/internal/validator"

	// Autoloads `.env`
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// Connect to database
	db, err := postgres.Open(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// Check if credentials are valid
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	// Create services
	as := postgres.NewAuthService(db)
	vs := validator.NewValidationService()
	cs := bcrypt.NewCryptoService()
	us := uuid.NewUUIDService()
	js := jwt.NewJWTService(os.Getenv("JWT_REFRESH_SECRET"), os.Getenv("JWT_ACCESS_SECRET"))
	r := chi.NewChiRouter()

	ac := chi.NewAuthController(as, vs, cs, us, js)

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
		r.Route("/", func(r chi.Router) {
			// r.Use(middleware.Authorization)
			r.Get("/", ac.GetUsers)
			r.Post("/", ac.CreateUser)
			r.Post("/login", ac.Login)
			r.Post("/refresh", ac.Refresh)
			r.Group(func(r chi.Router) {
				r.Use(chimiddleware.ValidateUUIDParam("UserUUID"))
				r.Get("/{UserUUID}", ac.GetUser)
				r.Patch("/{UserUUID}", ac.UpdateUser)
				r.Delete("/{UserUUID}", ac.DeleteUser)
			})
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
		log.Panic(err)
	}

}
