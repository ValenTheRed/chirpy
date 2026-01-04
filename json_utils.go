package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorPayload struct {
	Error string `json:"error"`
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
