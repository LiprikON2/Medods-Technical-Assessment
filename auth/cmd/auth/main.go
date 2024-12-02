package main

import (
	"log"

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
	h.ListenAndServe()
}
