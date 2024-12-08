package main

import (
	"log"
	"net/http"
	"os"
	"time"

	mddl "github.com/go-chi/chi/middleware"

	"github.com/medods-technical-assessment/internal/bcrypt"
	"github.com/medods-technical-assessment/internal/chi"
	cmddl "github.com/medods-technical-assessment/internal/chi/middleware"
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
	js := jwt.NewJWTService(os.Getenv("JWT_ACCESS_SECRET"), us)
	r := chi.NewChiRouter()

	ac := chi.NewAuthController(as, vs, cs, us, js)

	r.Use(mddl.StripSlashes)

	// A good base middleware stack
	r.Use(mddl.RequestID)
	// Not very trustworthy
	// ref: https://adam-p.ca/blog/2022/03/x-forwarded-for/#go-chichi
	r.Use(mddl.RealIP)
	r.Use(mddl.Logger)
	r.Use(mddl.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(mddl.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", ac.Register)
			r.Route("/login", func(r chi.Router) {
				r.With(cmddl.ValidateUUIDParam("UserUUID")).Post("/{UserUUID}", ac.LoginByUUID)
				r.Post("/", ac.Login)
			})
			r.Post("/refresh", ac.Refresh)

			r.With(cmddl.Authorization(js)).Get("/", ac.GetUsers)
			r.With(cmddl.Authorization(js)).Post("/", ac.CreateUser)

			r.With(cmddl.ValidateUUIDParam("UserUUID")).Group(func(r chi.Router) {
				r.Get("/{UserUUID}", ac.GetUser)
				r.With(cmddl.Authorization(js)).Patch("/{UserUUID}", ac.UpdateUser)
				r.With(cmddl.Authorization(js)).Delete("/{UserUUID}", ac.DeleteUser)
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Printf("Starting HTTP server at http://localhost:%s", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
