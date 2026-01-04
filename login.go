package main

import (
	"ValenTheRed/chirpy/internal/auth"
	"ValenTheRed/chirpy/internal/database"
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

	nullEmail := sql.NullString{
		Valid:  true,
		String: request.Email,
	}

	hashed_password, err := cfg.dbQueries.GetUsersHashedPassword(
		r.Context(),
		nullEmail,
	)
	if err != nil {
		log.Printf("Error in finding user's hashed_password: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	if match, err := auth.CheckPasswordHash(request.Password, hashed_password); err != nil {
		log.Printf("Error in comparing hashed_password and password: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	} else if !match {
		log.Printf("Error: users password does not match\n")
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	user, err := cfg.dbQueries.Login(r.Context(), database.LoginParams{
		Email:          nullEmail,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Printf("Error in login: %v\n", err)
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
