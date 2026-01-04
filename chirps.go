package main

import (
	"ValenTheRed/chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxChirpLength = 140
const profaneReplacement = "****"

var profanePattern = regexp.MustCompile(`(?i)(kerfuffle|sharbert|fornax)`)

func createChirpsHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body   string        `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
	}

	type responsePayload chirp

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

func listChirpsHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type responsePayloadItem chirp

	chirps, err := cfg.dbQueries.ListChirps(r.Context())
	if err != nil {
		log.Printf("Error when retrieving all chirps: %v\n", err)
		jsonResponse(w, http.StatusInternalServerError, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	response := make([]responsePayloadItem, 0, len(chirps))
	for _, chirp := range chirps {
		response = append(response, responsePayloadItem{
			ID:        chirp.ID,
			UserID:    chirp.UserID.UUID,
			Body:      chirp.Body.String,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
		})
	}
	jsonResponse(w, http.StatusOK, response)
}

func getChirpHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type responsePayload chirp

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error in finding chirp ID: %v\n", err)
		jsonResponse(w, http.StatusBadRequest, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error retrieving chirp of user: %v\n", err)
		// NOTE: ideally, should be matching error message and setting
		// status code basis that.
		jsonResponse(w, http.StatusNotFound, errorPayload{
			Error: "Something went wrong",
		})
		return
	}

	jsonResponse(w, http.StatusOK, responsePayload{
		ID:        chirp.ID,
		UserID:    chirp.UserID.UUID,
		Body:      chirp.Body.String,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
	})
}
