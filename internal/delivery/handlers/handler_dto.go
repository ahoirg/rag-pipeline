package handlers

import (
	"rag-pipeline/evaluation"
	"rag-pipeline/internal/services/rag_manager"
	"time"
)

type AskRequest struct {
	Query string `json:"query" validate:"required,min=3"`
}

type ApiResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Query     string    `json:"query,omitempty"`
	Answer    string    `json:"answer,omitempty"`
	Data      any       `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Handler struct {
	RagService *rag_manager.RAGService
	Evaluator  *evaluation.Evaluator
}
