package main

import (
	"ValenTheRed/chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	requestsCount atomic.Int64
	dbQueries     *database.Queries
	tokenSecret   string
}

func (cfg *apiConfig) increaseRequestsCount(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.requestsCount.Add(1)
		handler.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) logRequestsCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	response := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %v times!</p>
  </body>
</html>`, cfg.requestsCount.Load())
	w.Write(fmt.Append(nil, response))
}

func (cfg *apiConfig) resetRequestsCount(w http.ResponseWriter, r *http.Request) {
	if err := cfg.dbQueries.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("Error deleting all users: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cfg.requestsCount.Swap(0)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("could not connect to db at %v\n", dbUrl)
	}
	tokenSecret := os.Getenv("TOKEN_SECRET")

	cfg := apiConfig{
		dbQueries:   database.New(db),
		tokenSecret: tokenSecret,
	}
	root := os.DirFS(".")

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.increaseRequestsCount(http.StripPrefix("/app", http.FileServerFS(root))),
	)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", cfg.logRequestsCount)
	mux.HandleFunc("POST /admin/reset", enableOnDevEnv(cfg.resetRequestsCount))

	mux.HandleFunc("POST /api/login", withApiConfig(&cfg, loginHandler))
	mux.HandleFunc("POST /api/refresh", withApiConfig(&cfg, refreshHandler))
	mux.HandleFunc("POST /api/revoke", withApiConfig(&cfg, revokeHandler))

	mux.HandleFunc("POST /api/users", withApiConfig(&cfg, usersHandler))
	mux.HandleFunc("PUT /api/users", withApiConfig(&cfg, updateUsersHandler))

	mux.HandleFunc("POST /api/chirps", withApiConfig(&cfg, createChirpsHandler))
	mux.HandleFunc("GET /api/chirps", withApiConfig(&cfg, listChirpsHandler))
	mux.HandleFunc("GET /api/chirps/{chirpID}", withApiConfig(&cfg, getChirpHandler))
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", withApiConfig(&cfg, deleteChirpHandler))

	mux.HandleFunc("POST /api/polka/webhooks", withApiConfig(&cfg, polkaWebHooksHandler))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("server started at localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func enableOnDevEnv(handler http.HandlerFunc) http.HandlerFunc {
	if os.Getenv("PLATFORM") == "DEV" {
		return handler
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

func withApiConfig(
	cfg *apiConfig,
	handler func(cfg *apiConfig, w http.ResponseWriter, r *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(cfg, w, r)
	}
}
