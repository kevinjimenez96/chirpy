package types

type Chirp struct {
	Body string `json:"body"`
}

type AddChirpReq struct {
	Chirp
}

type ValidateChirpResponse struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}
