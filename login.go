package main

import (
	"ValenTheRed/chirpy/internal/auth"
	"ValenTheRed/chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const unauthorizedLoginErrorMessage = "Incorrect email or password"

func loginHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responsePayload struct {
		user
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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
		log.Printf("Error in comparing hashed password and password: %v\n", err)
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

	token, err := auth.MakeJWT(requester.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		log.Printf("Error in JWT token creation: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error when creating refresh token: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}
	if _, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: uuid.NullUUID{
			UUID:  requester.ID,
			Valid: true,
		},
	}); err != nil {
		log.Printf("Error when storing refresh token in database: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: unauthorizedLoginErrorMessage,
		})
		return
	}

	jsonResponse(w, http.StatusOK, responsePayload{
		user: user{
			ID:          requester.ID,
			CreatedAt:   requester.CreatedAt.Time,
			UpdatedAt:   requester.UpdatedAt.Time,
			Email:       requester.Email.String,
			IsChirpyRed: requester.IsChirpyRed.Bool,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func refreshHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type responsePayload struct {
		Token string `json:"token"`
	}

	headerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error in getting refresh token from authorization: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}
	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), headerRefreshToken)
	if err != nil {
		log.Printf("Erron in getting refresh token from database: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID.UUID, cfg.tokenSecret, time.Hour)
	if err != nil {
		log.Printf("Erron in creating JWT token: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	jsonResponse(w, http.StatusOK, responsePayload{
		Token: token,
	})
}

func revokeHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	authRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error in getting refresh token from authorization: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), authRefreshToken)
	if err != nil {
		log.Printf("Erron in revoking refresh token: %v\n", err)
		jsonResponse(w, http.StatusUnauthorized, errorPayload{
			Error: "Unauthorized",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
