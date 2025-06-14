package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantToken   string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Valid Bearer token",
			authHeader: "Bearer sometoken123",
			wantToken:  "sometoken123",
			wantErr:    false,
		},
		{
			name:        "No Authorization header",
			authHeader:  "",
			wantToken:   "",
			wantErr:     true,
			errContains: "No auth token",
		},
		{
			name:        "Invalid prefix",
			authHeader:  "Token sometoken123",
			wantToken:   "",
			wantErr:     true,
			errContains: "Invalid token",
		},
		{
			name:       "Bearer with no token",
			authHeader: "Bearer ",
			wantToken:  "",
			wantErr:    false,
		},
		{
			name:        "Bearer prefix in the middle",
			authHeader:  "sometoken Bearer abc",
			wantToken:   "",
			wantErr:     true,
			errContains: "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}
			gotToken, err := GetBearerToken(headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %q, want %q", gotToken, tt.wantToken)
			}
			if tt.wantErr && err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("GetBearerToken() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}
