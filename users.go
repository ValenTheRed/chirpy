package main

import (
	"ValenTheRed/chirpy/internal/auth"
	"ValenTheRed/chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func usersHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// NOTE: follows the generated model database.User
	type responsePayload user

	request := requestPayload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	hashed_password, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("Error hashing user's password: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: sql.NullString{
			Valid:  true,
			String: request.Email,
		},
		HashedPassword: hashed_password,
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
