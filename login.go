package main

import (
	"ValenTheRed/chirpy/internal/auth"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const unauthorizedLoginErrorMessage = "Incorrect email or password"

func loginHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email            string         `json:"email"`
		Password         string         `json:"password"`
		ExpiresInSeconds *time.Duration `json:"expires_in_seconds"`
	}

	type responsePayload struct {
		user
		Token string `json:"token"`
	}

	request := requestPayload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	requester, err := cfg.dbQueries.GetUser(r.Context(), sql.NullString{
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
		requester.HashedPassword,
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

	var expiresIn time.Duration
	if request.ExpiresInSeconds == nil || *request.ExpiresInSeconds > time.Hour {
		expiresIn = time.Hour
	} else {
		expiresIn = *request.ExpiresInSeconds
	}
	token, err := auth.MakeJWT(requester.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		log.Printf("Error in JWT token creation: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	jsonResponse(w, http.StatusOK, responsePayload{
		user: user{
			ID:        requester.ID,
			CreatedAt: requester.CreatedAt.Time,
			UpdatedAt: requester.UpdatedAt.Time,
			Email:     requester.Email.String,
		},
		Token: token,
	})
}
