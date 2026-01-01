package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const maxChirpLength = 140

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body string `json:"body"`
	}

	type responsePayload struct {
		Valid bool `json:"valid"`
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

	if len(request.Body) <= maxChirpLength {
		jsonResponse(w, http.StatusOK, responsePayload{
			Valid: true,
		})
	} else {
		jsonResponse(w, http.StatusBadRequest, errorPayload{
			Error: "Chirp is too long",
		})
	}
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
