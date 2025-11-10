package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"rag-pipeline/api"
)

func main() {

	// Initialize API dependencies
	api.InitApiDependencies()

	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Request logging
	r.Use(middleware.Recoverer) // Prevents server crash
	r.Mount("/", api.CreateRAGRouter())

	port := ":8080"
	// Start the server
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
