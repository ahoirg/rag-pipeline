package main

import (
	"log"
	"net/http"
	"rag-pipeline/configs"
	"rag-pipeline/internal/delivery/routers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("main.go|initialization error: %v", err)
	}

	// Initialize API dependencies
	sub_router, err := routers.CreateRAGRouter(config)
	if err != nil {
		log.Fatalf("main.go|initialization to start: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Request logging
	r.Use(middleware.Recoverer) // Prevents server crash
	r.Mount("/", sub_router)

	// Start the server
	port := ":" + config.Api.Port
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("main.go|Server failed to start: %v", err)
	}
}
