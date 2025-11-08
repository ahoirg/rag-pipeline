package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRAGRouter(t *testing.T) {

	// Initialize dependencies before creating router
	InitApiDependencies()

	// Create the router
	r := CreateRAGRouter()

	// GET /api/ping
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /api/ping: expected 200 OK, got %d", w.Code)
	}
}
