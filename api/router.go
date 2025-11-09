package api

import (
	"log"

	"github.com/go-chi/chi/v5"
)

// CreateRAGRouter creates and returns the API router
func CreateRAGRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/ping", PingHandler)
	r.Post("/api/ask", AskHandler)
	r.Post("/api/ask-directly", AskDirectlyHandler)
	r.Post("/api/ask-for-chunks", AskForChunksHandler)

	log.Println("   GET http://localhost:8080/api/ping")            // Health check endpoint
	log.Println("   POST http://localhost:8080/api/ask")            // Main RAG endpoint: question --> retrieval --> generation --> response
	log.Println("   POST http://localhost:8080/api/ask-directly")   //  question --> generation --> response
	log.Println("   POST http://localhost:8080/api/ask-for-chunks") // Endpoint to retrieve chunks relevant to the question

	return r
}

// InitApiDependencies initializes all services required by the API layer
func InitApiDependencies() error {
	return InitService()
}
