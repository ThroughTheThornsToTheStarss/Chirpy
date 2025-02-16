package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const (
	evt = "user.upgraded"
)

func (cfg *apiConfig) handlerSubscribe(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string            `json:"event"`
		Data  map[string]string `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Can't decode input parametes", err)
		return

	}
	if params.Event != evt {
		w.WriteHeader(204)
		return

	}

	userID, err := uuid.Parse(params.Data["user_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	err = cfg.db.RaiseUpUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "Can't find the user", err)
		return
	}
	w.WriteHeader(204)

}
