package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rag-pipeline/models"
	"time"
)

type OllamaEmbedder struct {
	BaseURL  string
	Endpoint string
	Model    string
	Client   *http.Client
}

// NewOllamaEmbedder creates and returns a new OllamaEmbedder
func NewOllamaEmbedder(baseUrl string, modelName string, endpoint string) *OllamaEmbedder {
	return &OllamaEmbedder{
		BaseURL:  baseUrl,
		Endpoint: endpoint,
		Model:    modelName,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// EmbedChunks sends the chunks to the Ollama embedder and returns their embeddings
func (e *OllamaEmbedder) EmbedChunks(chunks []string) ([][]float32, error) {

	reqBody := models.EmbedRequest{
		Model: e.Model,
		Input: chunks,
	}

	embedResp, err := e.embed(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to embed chunks: %w", err)
	}

	log.Printf(" Embedings size: %d", len(embedResp.Embeddings))

	return embedResp.Embeddings, nil
}

// EmbedQuery sends the query to the Ollama embedder and returns its embeddings
func (e *OllamaEmbedder) EmbedQuery(query string) ([]float32, error) {

	reqBody := models.EmbedRequest{
		Model: e.Model,
		Input: []string{query},
	}

	embedResp, err := e.embed(reqBody)
	if err != nil {
		return nil, fmt.Errorf("embeder.go|EmbedQuery: failed to embed the query: %w", err)
	} else if len(embedResp.Embeddings) <= 0 {
		return nil, fmt.Errorf("embeder.go|EmbedQuery: No embeddings found")
	}

	log.Printf("embeder.go|EmbedQuery: query embeding is completed")

	return embedResp.Embeddings[0], nil
}

// embed sends the given embedding request to the Ollama and returns the decoded embedding response
func (e *OllamaEmbedder) embed(reqBody models.EmbedRequest) (models.EmbedResponse, error) {

	var embedResp models.EmbedResponse

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return embedResp, err
	}

	resp, err := e.Client.Post(
		e.BaseURL+e.Endpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return embedResp, fmt.Errorf("embeder.go|embed: ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return embedResp, fmt.Errorf("embeder.go|embed: ollama returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return embedResp, fmt.Errorf("embeder.go|embed: failed to decode response: %w", err)
	}

	return embedResp, nil
}
