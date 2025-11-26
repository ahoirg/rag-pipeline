package routers

import (
	"log"
	"rag-pipeline/configs"
	"rag-pipeline/internal/delivery/handlers"

	"github.com/go-chi/chi/v5"
)

// CreateRAGRouter creates and returns the API router
func CreateRAGRouter(config *configs.Config) (*chi.Mux, error) {
	newHandler, err := handlers.NewHandler(config)
	if err != nil {
		return nil, err
	}
	r := chi.NewRouter()

	r.Get("/api/ping", handlers.PingHandler)
	r.Get("/api/evaluation", handlers.EvaluationHandler)
	r.Get("/api/evaluation/retrieval", newHandler.EvaluationRetrievalHandler)
	r.Get("/api/evaluation/generation", newHandler.EvaluationGenerationHandler)
	r.Post("/api/ask", newHandler.AskHandler)
	r.Post("/api/ask-directly", newHandler.AskDirectlyHandler)
	r.Post("/api/storebook", newHandler.StoreBookHandler)

	log.Println("   GET http://localhost:8080/api/ping")                  // Health check endpoint
	log.Println("   GET http://localhost:8080/api/evaluation/retrieval")  // get evaluation result of retrieval part
	log.Println("   GET http://localhost:8080/api/evaluation/generation") // get evaluation result of generation part
	log.Println("   POST http://localhost:8080/api/storebook")            // Store document into vector DB
	log.Println("   POST http://localhost:8080/api/ask")                  // Main RAG endpoint: question --> retrieval --> generation --> response
	log.Println("   POST http://localhost:8080/api/ask-directly")         // question --> generation --> response

	return r, nil
}
