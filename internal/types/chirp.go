package types

import "github.com/google/uuid"

type Chirp struct {
	Body string `json:"body"`
}

type AddChirpReq struct {
	Chirp
	UserId uuid.UUID `json:"user_id"`
}


type ValidateChirpResponse struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}
