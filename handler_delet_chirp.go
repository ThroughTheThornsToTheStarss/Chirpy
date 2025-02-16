package main

import (
	"chirpy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeletChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "Can't find token", err)
		return
	}
	UserID, err := auth.ValidateJWT(token, cfg.tS)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't bring userID", err)
		return
	}

	if dbChirp.UserID == UserID {
		err = cfg.db.DeleteChirp(r.Context(), chirpID)
		if err != nil {
			respondWithError(w, 403, "Failed to delete chirp", err)
		}
		w.WriteHeader(204)
		return

	} else {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp", nil)
	}
}
