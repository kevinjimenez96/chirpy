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

	expiresInSeconds := 60 * 60
	if loginUserReq.ExpiresInSeconds > 0 {
		expiresInSeconds = loginUserReq.ExpiresInSeconds
	}

	token, err := auth.MakeJWT(loggedUser.ID, cfg.Secret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error: token creation error: %s", err), err)
		return
	}
	respondWithJSON(w, http.StatusOK, types.LoginUserRes{
		ID:        loggedUser.ID,
		CreatedAt: loggedUser.CreatedAt,
		UpdatedAt: loggedUser.UpdatedAt,
		Email:     loggedUser.Email,
		Token:     token,
	})
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
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	})
}
