package handlers

import (
	"net/http"

	"github.com/kevinjimenez96/chirpy/internal/types"
)

func ResetHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
	}
	err := cfg.DbQueries.DeleteAllUsers(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error: could not reset users", err)
	}
	cfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
