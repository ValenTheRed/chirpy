package main

import (
	"ValenTheRed/chirpy/internal/auth"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

const unauthorizedLoginErrorMessage = "Incorrect email or password"

func loginHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responsePayload user

	request := requestPayload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	user, err := cfg.dbQueries.GetUser(r.Context(), sql.NullString{
		Valid:  true,
		String: request.Email,
	})
	if err != nil {
		log.Printf("Error in finding user: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	if match, err := auth.CheckPasswordHash(
		request.Password,
		user.HashedPassword,
	); err != nil {
		log.Printf("Error in comparing hashed_password and password: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	} else if !match {
		log.Printf("Error: user's password does not match\n")
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	jsonResponse(w, http.StatusOK, responsePayload{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email.String,
	})
}
