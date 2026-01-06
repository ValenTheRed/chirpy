package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type polkaEvent string

const userUpgradedEvent polkaEvent = "user.upgraded"

func polkaWebHooksHandler(cfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Event polkaEvent `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	request := requestPayload{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("POST polka webhooks: error in parsing request body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if request.Event != userUpgradedEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err := cfg.dbQueries.UpgradeUserToRed(r.Context(), request.Data.UserID)
	if err != nil {
		log.Printf("POST polka webhooks: error in upgrading user to Chirpy Red: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
