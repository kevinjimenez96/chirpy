package types

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoginUserRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}
