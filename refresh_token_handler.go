package main

import (
	"net/http"
	"time"

	"github.com/Dorfieeee/bootdev-http-server/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbRefreshTokenUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenString)
	if err != nil || dbRefreshTokenUser.RevokedAt.Valid || dbRefreshTokenUser.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
		return
	}

	token, err := auth.MakeJWT(dbRefreshTokenUser.ID, cfg.appSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new token", err)
		return
	}

	type RefreshTokenResponse struct {
		AccessToken string `json:"token"`
	}

	respondWithJSON(
		w, http.StatusOK, RefreshTokenResponse{
			AccessToken: token,
		},
	)
}
