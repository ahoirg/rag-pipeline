package api

import (
	"net/http"
)

// InitService initializes the api service
func InitService() {
	// TODO: chunker,embeddings,rag, database ?
	// or only a single service that encapsulates all?
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
