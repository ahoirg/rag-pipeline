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
	log.Println("   GET http://localhost:8080/api/ping") // Health check endpoint
	//log.Println("   POST http://localhost:8080/api/ask")          // Main RAG endpoint: question --> retrieval --> generation --> response
	//log.Println("   POST http://localhost:8080/api/askforchunks") // Endpoint to retrieve chunks relevant to the question

	// Start the server
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
