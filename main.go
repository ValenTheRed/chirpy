package main

import (
	"net/http"
	"os"
)

func main() {
	root := os.DirFS(".")

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServerFS(root)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
