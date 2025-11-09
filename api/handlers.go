package api

import (
	"net/http"
	"rag-pipeline/services"
)

var ragService *services.RAGService

// InitService initializes the api service
func InitService() error {
	var err error
	ragService, err = services.NewRAGService()
	return err
}

// PingHandler handles the health check endpoint
func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "RAG Pipeline API is running",
	}

	writeJSON(w, http.StatusOK, response)
}

// TODO
func AskHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}
