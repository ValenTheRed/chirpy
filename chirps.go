package main

import (
	"ValenTheRed/chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxChirpLength = 140
const profaneReplacement = "****"

var profanePattern = regexp.MustCompile(`(?i)(kerfuffle|sharbert|fornax)`)

func chirpsHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body   string        `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
	}

	type responsePayload struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
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

	if utf8.RuneCountInString(request.Body) > maxChirpLength {
		jsonResponse(w, http.StatusBadRequest, errorPayload{
			Error: "Chirp is too long",
		})
		return
	}

	cleanedBody := profanePattern.ReplaceAllString(request.Body, profaneReplacement)
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: sql.NullString{
			Valid:  true,
			String: cleanedBody,
		},
		UserID: request.UserID,
	})
	if err != nil {
		log.Printf("Error writing chirp: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	jsonResponse(w, http.StatusCreated, responsePayload{
		ID:        chirp.ID,
		UserID:    chirp.UserID.UUID,
		Body:      chirp.Body.String,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
	})
}
