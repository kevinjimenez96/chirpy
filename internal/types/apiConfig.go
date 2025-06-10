package types

import (
	"net/http"
	"sync/atomic"

	"github.com/kevinjimenez96/chirpy/internal/auth"
	"github.com/kevinjimenez96/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	Secret         string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareAddConfig(handler func(http.ResponseWriter, *http.Request, *ApiConfig)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, cfg)
	})
}

func (cfg *ApiConfig) MiddlewareAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err = auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
