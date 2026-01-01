package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	requestsCount atomic.Int64
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
	cfg.requestsCount.Swap(0)
	w.Write([]byte("OK"))
}

func main() {
	cfg := apiConfig{}
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
	mux.HandleFunc("POST /admin/reset", cfg.resetRequestsCount)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("server started at localhost:8080")
	log.Fatal(server.ListenAndServe())
}
