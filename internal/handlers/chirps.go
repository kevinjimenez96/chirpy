package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/kevinjimenez96/chirpy/internal/auth"
	"github.com/kevinjimenez96/chirpy/internal/database"
	"github.com/kevinjimenez96/chirpy/internal/types"
)

func GetAllChirps(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	chirps, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting all chirps: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func GetChirpById(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error: %s", err), err)
		return
	}

	chirp, err := cfg.DbQueries.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting chirp: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func AddChirp(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	// this has been already checked
	token, _ := auth.GetBearerToken(r.Header)
	userId, _ := auth.ValidateJWT(token, cfg.Secret)

	decoder := json.NewDecoder(r.Body)
	addChirp := types.AddChirpReq{}
	err := decoder.Decode(&addChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding chirp: %s", err), err)
		return
	}

	if len(addChirp.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp too short: %d chars long", len(addChirp.Body)), nil)
		return
	}

	if len(addChirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp too long: %d chars long", len(addChirp.Body)), nil)
		return
	}

	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userId,
		Body:   censorChrip(addChirp.Body),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error saving chirp: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func censorChrip(chirp string) string {
	censoredWords := getCensoredWords()

	cleanedChirp := strings.Split(chirp, " ")

	for i, word := range cleanedChirp {
		if censoredWords[strings.ToLower(word)] {
			cleanedChirp[i] = "****"
		}
	}

	return strings.Join(cleanedChirp, " ")
}

func getCensoredWords() map[string]bool {
	censoredWords := make(map[string]bool)
	censoredWords["kerfuffle"] = true
	censoredWords["sharbert"] = true
	censoredWords["fornax"] = true
	return censoredWords
}
