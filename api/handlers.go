package api

import (
	"encoding/json"
	"io"
	"net/http"
	"rag-pipeline/evaluation"
	"rag-pipeline/models"
	"rag-pipeline/services"
	"time"
)

var ragService *services.RAGService
var evaluator *evaluation.Evaluator

// InitService initializes the api service
func InitService(config *models.Config) error {
	var err error
	ragService, err = services.NewRAGService(config, config.Api.CollectionName)
	if err != nil {
		return err
	}

	evaluator, err = evaluation.NewEvaluator(config)
	if err != nil {
		return err
	}

	return nil
}

// PingHandler handles the health check endpoint
func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := models.ApiResponse{
		Success:   true,
		Message:   "Api is working...",
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}

// AskHandler handles the ask endpoint
func AskHandler(w http.ResponseWriter, r *http.Request) {
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

	response := models.ApiResponse{
		Success:   true,
		Query:     req.Query,
		Answer:    generatedResponse,
		Data:      chunks,
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}

// AskHandler handles the ask-directly endpoint
// bypasses RAG and uses only the LLM to generate a response
func AskDirectlyHandler(w http.ResponseWriter, r *http.Request) {
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

	response := models.ApiResponse{
		Success:   true,
		Query:     req.Query,
		Answer:    generatedResponse,
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}

// StoreBookHandler is endpoint to store document into vector DB
func StoreBookHandler(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Set the key: 'file' ", err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "File could not be converted to text. Send only .txt type files ", err)
		return
	}

	err = ragService.StoreData(string(content))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to store teh data", err)
		return
	}

	response := models.ApiResponse{
		Success:   true,
		Message:   "File data stored successfully!",
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}

// EvaluationGenerationHandler returns the evaluation results
// of the generation part of the RAGpipeline with the eval data
func EvaluationGenerationHandler(w http.ResponseWriter, r *http.Request) {
	result, err := evaluator.GetGenerationEvaluateResult()

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Generator evaluation in Rag pipeline could not be done: ", err)
		return
	}

	response := models.ApiResponse{
		Success:   true,
		Data:      result,
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}

// EvaluationRetrievalHandler returns the evaluation results
// of the Retrieval part of the RAGpipeline with the eval data
func EvaluationRetrievalHandler(w http.ResponseWriter, r *http.Request) {
	result, err := evaluator.GetRetrievalEvaluateResult()

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Retrieval evaluation in Rag pipeline could not be done: ", err)
		return
	}

	response := models.ApiResponse{
		Success:   true,
		Data:      result,
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}

// EvaluationHandler only returs Retrieval and Generation endpoints
func EvaluationHandler(w http.ResponseWriter, r *http.Request) {
	message := "For Retrieval Evaluation:  GET http://localhost:8080/api/evaluation/retrieval"
	message += "-- For Generation Evaluation:  GET http://localhost:8080/api/evaluation/generation"

	response := models.ApiResponse{
		Success:   true,
		Message:   message,
		Timestamp: time.Now(),
	}

	writeJSON(w, http.StatusOK, response)
}
