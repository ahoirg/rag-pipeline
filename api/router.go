package api

import (
	"log"
	"rag-pipeline/models"

	"github.com/go-chi/chi/v5"
)

// CreateRAGRouter creates and returns the API router
func CreateRAGRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/ping", PingHandler)
	r.Get("/api/evaluation", EvaluationHandler)
	r.Get("/api/evaluation/retrieval", EvaluationRetrievalHandler)
	r.Get("/api/evaluation/generation", EvaluationGenerationHandler)
	r.Post("/api/ask", AskHandler)
	r.Post("/api/ask-directly", AskDirectlyHandler)
	r.Post("/api/storebook", StoreBookHandler)

	log.Println("   GET http://localhost:8080/api/ping")                  // Health check endpoint
	log.Println("   GET http://localhost:8080/api/evaluation/retrieval")  // get evaluation result of retrieval part
	log.Println("   GET http://localhost:8080/api/evaluation/generation") // get evaluation result of generation part
	log.Println("   POST http://localhost:8080/api/storebook")            // Store document into vector DB
	log.Println("   POST http://localhost:8080/api/ask")                  // Main RAG endpoint: question --> retrieval --> generation --> response
	log.Println("   POST http://localhost:8080/api/ask-directly")         // question --> generation --> response

	return r
}

// InitApiDependencies initializes all services required by the API layer
func InitApiDependencies(config *models.Config) {
	InitService(config)
}
