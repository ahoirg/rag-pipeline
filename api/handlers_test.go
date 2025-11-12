package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {

	// Create new Http request
	req := httptest.NewRequest("GET", "/api/ping", nil)

	// Create ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the PingHandler
	PingHandler(w, req)

	// Check if the API request was successful
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check if the response Content-Type header is JSON
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Decode the JSON response body into a map
	var response map[string]any
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	//Check that the 'status' field value is 'success'
	status, ok := response["status"]
	if !ok {
		t.Error("Response missing 'status' field")
	}
	if status != "success" {
		t.Errorf("Expected status 'success', got '%v'", status)
	}
}

func TestAskHandler(t *testing.T) {
	// TODO
}
