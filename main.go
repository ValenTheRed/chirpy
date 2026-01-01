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

func (cfg *apiConfig) logRequestsCount() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fmt.Appendf(nil, "Hits: %v", cfg.requestsCount.Load()))
	})
}

func (cfg *apiConfig) resetRequestsCount() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.requestsCount.Swap(0)
		w.Write([]byte("OK"))
	})
}

func main() {
	cfg := apiConfig{}
	root := os.DirFS(".")

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.increaseRequestsCount(http.StripPrefix("/app", http.FileServerFS(root))),
	)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("GET /metrics", cfg.logRequestsCount())
	mux.Handle("POST /reset", cfg.resetRequestsCount())

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("server started at localhost:8080")
	log.Fatal(server.ListenAndServe())
}
