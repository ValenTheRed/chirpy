package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"unicode/utf8"
)

const maxChirpLength = 140
const profaneReplacement = "****"

var profanePattern = regexp.MustCompile(`(?i)(kerfuffle|sharbert|fornax)`)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body string `json:"body"`
	}

	type responsePayload struct {
		CleanedBody string `json:"cleaned_body"`
	}

	type errorPayload struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)

	request := requestPayload{}
	err := decoder.Decode(&request)
	if err != nil {
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
	jsonResponse(w, http.StatusOK, responsePayload{
		CleanedBody: cleanedBody,
	})
}

func jsonResponse(w http.ResponseWriter, status int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(status)
	}
	w.Header().Set("Content-Type", "text/json")
	w.Write(response)
}
