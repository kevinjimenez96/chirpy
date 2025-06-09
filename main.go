package main

import (
	"log"
	"net/http"

	"github.com/kevinjimenez96/chirpy/internal/handlers"
	"github.com/kevinjimenez96/chirpy/internal/types"
)

func main() {
	var cfg = &types.ApiConfig{}

	port := "8080"
	filepathRoot := http.Dir(".")

	serveMux := http.NewServeMux()

	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(filepathRoot))))

	serveMux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", handlers.ValidateChirpHandler)

	serveMux.Handle("GET /admin/metrics", cfg.MiddlewareAddConfig(handlers.MetricsHandler))
	serveMux.Handle("POST /admin/reset", cfg.MiddlewareAddConfig(handlers.ResetHandler))

	srv := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
