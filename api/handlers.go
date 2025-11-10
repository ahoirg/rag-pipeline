package api

import (
	"encoding/json"
	"io"
	"net/http"
	"rag-pipeline/evaluation"
	"rag-pipeline/models"
	"rag-pipeline/services"
)

var ragService *services.RAGService
var evaluator *evaluation.Evaluator

const apiCollectionName = "api_collection"

// InitService initializes the api service
func InitService() {
	ragService = services.NewRAGService(apiCollectionName)
	evaluator = evaluation.NewEvaluator()
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

	generatedResponse, chunks, err := ragService.GenerateResponse(req.Query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate the response", err)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": generatedResponse,
		"data":    chunks,
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

// StoreBookHandler is endpoint to store document into vector DB
func StoreBookHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusMethodNotAllowed, "Set the key: 'file' ", err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "File could not be converted to text", err)
		return
	}

	err = ragService.StoreData(string(content))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to store data", err)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "File data stored successfully",
	}

	writeJSON(w, http.StatusOK, response)
}

// EvaluationGenerationHandler returns the evaluation results
// of the generation part of the RAGpipeline with the eval data
func EvaluationGenerationHandler(w http.ResponseWriter, r *http.Request) {
	result, err := evaluator.GetGenerationEvaluateResult()

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Evaluation ERROR", err)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"data":   result,
	}

	writeJSON(w, http.StatusOK, response)
}

// EvaluationRetrievalHandler returns the evaluation results
// of the Retrieval part of the RAGpipeline with the eval data
func EvaluationRetrievalHandler(w http.ResponseWriter, r *http.Request) {
	result, err := evaluator.GetRetrievalEvaluateResult()

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Evaluation ERROR", err)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"data":   result,
	}

	writeJSON(w, http.StatusOK, response)
}
