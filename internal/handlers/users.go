package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kevinjimenez96/chirpy/internal/auth"
	"github.com/kevinjimenez96/chirpy/internal/database"
	"github.com/kevinjimenez96/chirpy/internal/types"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	loginUserReq := types.LoginUserReq{}
	err := decoder.Decode(&loginUserReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding login user request: %s", err), err)
		return
	}

	loggedUser, err := cfg.DbQueries.GetUserByEmail(r.Context(), loginUserReq.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error looking for user: %s", err), err)
		return
	}
	err = auth.CheckPasswordHash(loginUserReq.Password, loggedUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error: Incorrect email or passwordt: %s", err), err)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	cfg.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID: loggedUser.ID,
		Token:  refreshToken,
	})

	token, err := auth.MakeJWT(loggedUser.ID, cfg.Secret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error: token creation error: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusOK, types.LoginUserRes{
		ID:           loggedUser.ID,
		CreatedAt:    loggedUser.CreatedAt,
		UpdatedAt:    loggedUser.UpdatedAt,
		Email:        loggedUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  loggedUser.IsChirpyRed.Bool,
	})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error invalid token: %s", err), err)
		return
	}

	refreshToken, err := cfg.DbQueries.GetRefreshToken(r.Context(), token)
	if err != nil || !time.Time.IsZero(refreshToken.RevokedAt.Time) {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error looking for token or expired: %s", err), err)
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.Secret, time.Duration(1)*time.Hour)
	if err != nil || time.Time.IsZero(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating new token: %s", err), err)
		return
	}
	respondWithJSON(w, http.StatusOK, types.RefreshTokenRes{
		Token: newToken,
	})
}

func RevokeHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error invalid token: %s", err), err)
		return
	}

	_, err = cfg.DbQueries.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error revoking for token: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func AddUserHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	decoder := json.NewDecoder(r.Body)
	createUserReq := types.CreateUserReq{}
	err := decoder.Decode(&createUserReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding create user request: %s", err), err)
		return
	}

	hashedPassword, err := auth.HashPassword(createUserReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %s", err), err)
		return
	}

	newUser, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          createUserReq.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding creating user: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, types.LoginUserRes{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed.Bool,
	})
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	token, _ := auth.GetBearerToken(r.Header)
	userId, _ := auth.ValidateJWT(token, cfg.Secret)

	decoder := json.NewDecoder(r.Body)
	updateUserReq := types.UpdateUserReq{}
	err := decoder.Decode(&updateUserReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding update user request: %s", err), err)
		return
	}

	hashedPassword, err := auth.HashPassword(updateUserReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %s", err), err)
		return
	}

	newUser, err := cfg.DbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          updateUserReq.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding updating user: %s", err), err)
		return
	}

	respondWithJSON(w, http.StatusOK, types.LoginUserRes{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed.Bool,
	})
}
