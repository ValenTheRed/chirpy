package main

import (
	"net/http"
	"os"
)

func main() {
	root := os.DirFS(".")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServerFS(root))
	mux.Handle("/assets/logo.png", http.FileServerFS(root))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
