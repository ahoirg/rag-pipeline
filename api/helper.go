package api

import (
	"encoding/json"
	"net/http"
	"rag-pipeline/models"
	"time"
)

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError is a helper function to write error responses
func writeError(w http.ResponseWriter, status int, message string, err error) {

	response := models.ApiResponse{
		Success:   false,
		Message:   message + err.Error(),
		Timestamp: time.Now(),
	}

	writeJSON(w, status, response)
}
