package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kevinjimenez96/chirpy/internal/types"
)

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := types.Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding chirp: %s", err), err)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp too long: %d chars long", len(chirp.Body)), nil)
		return
	}

	validateChirpResponse := types.ValidateChirpResponse{
		Valid:       true,
		CleanedBody: censorChrip(chirp.Body),
	}

	respondWithJSON(w, http.StatusOK, validateChirpResponse)
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
