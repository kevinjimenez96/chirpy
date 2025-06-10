package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinjimenez96/chirpy/internal/auth"
	"github.com/kevinjimenez96/chirpy/internal/database"
	"github.com/kevinjimenez96/chirpy/internal/types"
)

func PolkaWebHook(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	key, err := auth.GetAPIKey(r.Header)
	if key != cfg.PolkaKey || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: not authorize", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	polkaWebHookReq := types.PolkaWebHookReq{}
	err = decoder.Decode(&polkaWebHookReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding polka webhook request: %s", err), err)
		return
	}

	if polkaWebHookReq.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.DbQueries.UpdateUserIsChirpyRedById(r.Context(), database.UpdateUserIsChirpyRedByIdParams{
		ID: polkaWebHookReq.Data.UserId,
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error updating user: polka webhook request: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
