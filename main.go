package main

import (
	"fmt"
	"log"
	"net/http"

	"rag-pipeline/api"
	"rag-pipeline/models"
	"rag-pipeline/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/yaml.v3"
)

const config_path = "config.yaml"

func main() {

	var config models.Config
	if err := loadConfig(&config); err != nil {
		log.Fatalf("main.go|initialization error: %v", err)
	}

	// Initialize API dependencies
	api.InitApiDependencies(&config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Request logging
	r.Use(middleware.Recoverer) // Prevents server crash
	r.Mount("/", api.CreateRAGRouter())

	// Start the server
	port := ":" + config.Api.Port
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("main.go|Server failed to start: %v", err)
	}
}

func loadConfig(config *models.Config) error {
	raw_config, err := utils.LoadDocument(config_path)
	if err != nil {
		return fmt.Errorf("main.go|failed to load config.yaml: %w", err)
	}

	if err := yaml.Unmarshal([]byte(raw_config), &config); err != nil {
		return fmt.Errorf("main.go|failed to convert raw_config to yaml: %w", err)
	}

	return nil
}
