package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/medods-technical-assessment"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Print("GetUser")

	response := auth.User{
		ID:       123,
		Email:    "test@example.com",
		Password: "XXXXXXXX",
		// Code: http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
		// auth.InternalErrorHandler(w)
		return
	}

}
