package main

import (
	"chirpy/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusOK, Chirps{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerReturn(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error

	sortParam := r.URL.Query().Get("sort")
	authorID := r.URL.Query().Get("author_id")

	if authorID != "" {
		chirpID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
			return
		}
		chirps, err = cfg.db.GetChirpByUser(r.Context(), chirpID)
	} else if sortParam == "desc" {
		chirps, err = cfg.db.GetChirpsDECS(r.Context())
	} else {
		// Default case (asc or no sort parameter)
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't return parameters", err)
		return
	}

	res := getmiddlchirp(chirps)
	respondWithJSON(w, http.StatusOK, res)
}

func getmiddlchirp(dbChirps []database.Chirp) []Chirps {
	chirpsDtos := make([]Chirps, 0, len(dbChirps))

	for i := 0; i < len(dbChirps); i++ {
		chirpsDtos = append(chirpsDtos, Chirps{
			ID:        dbChirps[i].ID,
			CreatedAt: dbChirps[i].CreatedAt,
			UpdatedAt: dbChirps[i].UpdatedAt,
			Body:      dbChirps[i].Body,
			UserID:    dbChirps[i].UserID,
		})
	}
	return chirpsDtos
}
