package types

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type RefreshTokenRes struct {
	Token string `json:"token"`
}
