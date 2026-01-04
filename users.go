package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func usersHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email string `json:"email"`
	}

	// NOTE: follows the generated model database.User
	type responsePayload struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	type errorPayload struct {
		Error string `json:"error"`
	}

	request := requestPayload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), sql.NullString{
		Valid:  true,
		String: request.Email,
	})
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	jsonResponse(w, http.StatusCreated, responsePayload{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email.String,
	})
}
