package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError is a helper function to write error responses
func writeError(w http.ResponseWriter, status int, message string, err error) {

	response := map[string]interface{}{
		"status":  "error",
		"message": message + ": " + err.Error(),
	}

	writeJSON(w, status, response)
}
