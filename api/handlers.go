package api

import (
	"encoding/json"
	"net/http"
	"rag-pipeline/models"
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
		"data":    nil,
	}

	writeJSON(w, http.StatusOK, response)
}

// AskHandler handles the ask endpoint
func AskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req models.AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	generatedResponse, err := ragService.GenerateResponse(req.Query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate the response", err)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": generatedResponse,
	}

	writeJSON(w, http.StatusOK, response)
}

// AskHandler handles the ask-directly endpoint
// bypasses RAG and uses only the LLM to generate a response
func AskDirectlyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req models.AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	generatedResponse, err := ragService.GenerateResponseWithoutChunks(req.Query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate the response", err)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": generatedResponse,
	}

	writeJSON(w, http.StatusOK, response)
}

// AskDirectlyHandler handles the ask-for-chunks endpoint
// retrieves relevant chunks without generating a full response
func AskForChunksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req models.AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	chunks, err := ragService.RetrieveRelevantChunks(req.Query, 5)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve relevant chunks", err)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "RAG Pipeline API is running",
		"data":    chunks,
	}

	writeJSON(w, http.StatusOK, response)
}
