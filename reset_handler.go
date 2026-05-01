package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		return
	}

	cfg.fileserverHits.Store(0)
	if err := cfg.db.DeleteAllUsers(req.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete all users", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset:utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
