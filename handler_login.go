package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"time"
)

const (
	maxExpirationTime = 3600
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > maxExpirationTime {
		params.ExpiresInSeconds = maxExpirationTime
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't found user", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong email or passord", err)
	}
	expiresTime := time.Duration(params.ExpiresInSeconds) * time.Second
	token, err := auth.MakeJWT(user.ID, cfg.tS, expiresTime)
	if err != nil {
		respondWithError(w, 401, "Couldn't make JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			token:     token,
		},
		Token: token,
	})
}
