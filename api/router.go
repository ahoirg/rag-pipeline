package api

import (
	"github.com/go-chi/chi/v5"
)

// CreateRAGRouter creates and returns the API router
func CreateRAGRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/ping", PingHandler)
	r.Post("/api/ask", AskHandler)

	return r
}

// InitApiDependencies initializes all services required by the API layer
func InitApiDependencies() {
	InitService()
}
