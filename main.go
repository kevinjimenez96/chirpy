package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kevinjimenez96/chirpy/internal/database"
	"github.com/kevinjimenez96/chirpy/internal/handlers"
	"github.com/kevinjimenez96/chirpy/internal/types"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Error opening db connection.")
	}

	var cfg = &types.ApiConfig{
		DbQueries: database.New(db),
		Platform:  os.Getenv("PLATFORM"),
		Secret:    os.Getenv("SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}

	port := "8080"
	filepathRoot := http.Dir(".")

	serveMux := http.NewServeMux()

	serveMux.Handle("GET /admin/metrics", cfg.MiddlewareAddConfig(handlers.MetricsHandler))
	serveMux.Handle("POST /admin/reset", cfg.MiddlewareAddConfig(handlers.ResetHandler))
	serveMux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)

	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(filepathRoot))))

	serveMux.Handle("GET /api/chirps", cfg.MiddlewareAddConfig(handlers.GetAllChirps))
	serveMux.Handle("GET /api/chirps/{id}", cfg.MiddlewareAddConfig(handlers.GetChirpById))
	serveMux.Handle("POST /api/chirps", cfg.MiddlewareAuth(cfg.MiddlewareAddConfig(handlers.AddChirp)))
	serveMux.Handle("DELETE /api/chirps/{id}", cfg.MiddlewareAuth(cfg.MiddlewareAddConfig(handlers.DeleteChirpByIdHandler)))

	serveMux.Handle("POST /api/users", cfg.MiddlewareAddConfig(handlers.AddUserHandler))
	serveMux.Handle("PUT /api/users", cfg.MiddlewareAuth(cfg.MiddlewareAddConfig(handlers.UpdateUserHandler)))

	serveMux.Handle("POST /api/login", cfg.MiddlewareAddConfig(handlers.LoginHandler))
	serveMux.Handle("POST /api/refresh", cfg.MiddlewareAddConfig(handlers.RefreshTokenHandler))
	serveMux.Handle("POST /api/revoke", cfg.MiddlewareAddConfig(handlers.RevokeHandler))

	serveMux.Handle("POST /api/polka/webhooks", cfg.MiddlewareAddConfig(handlers.PolkaWebHook))

	srv := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
