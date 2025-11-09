package models

import "net/http"

type OllamaEmbedderConfig struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

type EmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbedResponse struct {
	Embedding  []float32   `json:"embedding"`
	Embeddings [][]float32 `json:"embeddings"`
}
