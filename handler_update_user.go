package main

import (
	"encoding/json"
	"net/http"

	"github.com/svladivanov/chirpy/internal/auth"
	"github.com/svladivanov/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't verify bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate bearer token", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	updatedDBUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPw,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	apiUser := User{
		ID:          updatedDBUser.ID,
		CreatedAt:   updatedDBUser.CreatedAt,
		UpdatedAt:   updatedDBUser.UpdatedAt,
		Email:       updatedDBUser.Email,
		IsChirpyRed: updatedDBUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, apiUser)
}
