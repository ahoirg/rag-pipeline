package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"rag-pipeline/models"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCreateRAGRouter(t *testing.T) {
	path := filepath.Clean("../config.yaml") //TODO

	raw_config, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("TestCreateRAGRouter FAIL: readfile, detail: %s", err.Error())
	}

	var config models.Config
	if err := yaml.Unmarshal([]byte(raw_config), &config); err != nil {
		t.Errorf("TestCreateRAGRouter FAIL: failed to convert raw_config to yaml, detail: %s", err.Error())
	}

	// Initialize dependencies before creating router
	if err := InitApiDependencies(&config); err != nil {
		t.Errorf("GET /api/ping: expected 200 OK, got %s", err.Error())
	}

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
