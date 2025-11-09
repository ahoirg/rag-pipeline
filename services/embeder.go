package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rag-pipeline/models"
)

// EmbedChunks sends the chunks to the Ollama embedder and returns their embeddings
func EmbedChunks(chunks []string, e models.OllamaEmbedderConfig) ([][]float32, error) {

	reqBody := models.EmbedRequest{
		Model: e.Model,
		Input: chunks,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := e.Client.Post(
		e.BaseURL+"/api/embed",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var embedResp models.EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf(" Embedings size: %d", len(embedResp.Embeddings))

	return embedResp.Embeddings, nil
}
