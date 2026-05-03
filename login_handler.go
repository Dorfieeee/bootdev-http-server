package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dorfieeee/bootdev-http-server/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type loginReqParams struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	defaultTokenDuration := time.Duration(time.Hour)
	decoder := json.NewDecoder(r.Body)
	var params loginReqParams
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Incorrect email or password", err)
		return
	}

	isAuth, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Authentication failed", err)
		return
	}

	if !isAuth {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiresIn := defaultTokenDuration
	if params.ExpiresInSeconds != nil {
		userDefinedDuration := time.Duration(*params.ExpiresInSeconds) * time.Second
		if userDefinedDuration < defaultTokenDuration {
			expiresIn = userDefinedDuration
		}
	}

	JWTToken, err := auth.MakeJWT(user.ID, cfg.appSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}

	type loginResponse struct {
		User
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.CreatedAt,
			Email:     user.Email,
		},
		Token: JWTToken,
	})
}
