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

	hashedPassword, err := auth.HashPassword(request.Password)
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
		HashedPassword: hashedPassword,
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

func updateUsersHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type responsePayload user

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error in getting authorization token from header: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}
	request := requestPayload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("Error in hashing the password: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		log.Printf("Error validating JWT token: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: sql.NullString{
			String: request.Email,
			Valid:  true,
		},
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error updating user's email and password: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	jsonResponse(w, http.StatusOK, responsePayload{
		ID:        userID,
		Email:     request.Email,
		UpdatedAt: user.UpdatedAt.Time,
		CreatedAt: user.CreatedAt.Time,
	})
}
