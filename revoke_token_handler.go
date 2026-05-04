package main

import (
	"net/http"

	"github.com/Dorfieeee/bootdev-http-server/internal/auth"
)

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), refreshTokenString); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to revoke token: "+err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
