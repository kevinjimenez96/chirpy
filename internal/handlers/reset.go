package handlers

import (
	"net/http"

	"github.com/kevinjimenez96/chirpy/internal/types"
)

func ResetHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	cfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
