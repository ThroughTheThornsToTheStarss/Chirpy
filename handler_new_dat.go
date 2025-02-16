package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

type parameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type response struct {
	User
}

func (cfg *apiConfig) handlerNewdat(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "Can't fined token", err)
		return
	}
	UserID, err := auth.ValidateJWT(token, cfg.tS)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't bring userID", err)
		return
	}

	hashPsrd, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't hash passord", err)
	}
	updateUser, err := cfg.db.Update(r.Context(), database.UpdateParams{
		Email:          param.Email,
		HashedPassword: hashPsrd,
		ID:             UserID,
	})
	if err != nil {
		respondWithError(w, 500, "Can't update the users dat", err)
	}
	respondWithJSON(w, 200, response{
		User: User{
			ID:          updateUser.ID,
			CreatedAt:   updateUser.UpdatedAt,
			UpdatedAt:   updateUser.UpdatedAt,
			Email:       updateUser.Email,
			IsChirpyRed: updateUser.IsChirpyRed,
		},
	})

}
