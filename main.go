package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"slices"
	"sync/atomic"

	"github.com/Dorfieeee/bootdev-http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const filepathRoot string = "."
const port string = "8080"

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	appSecret      string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	appSecret := os.Getenv("APP_SECRET")
	platform := os.Getenv("PLATFORM")

	if slices.Contains([]string{dbURL, appSecret, platform}, "") {
		log.Fatal("Env variables are not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to open db connection: %v\n", err.Error())
	}
	dbQueries := database.New(db)

	cfg := apiConfig{
		db:        dbQueries,
		platform:  platform,
		appSecret: appSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
