package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LLMService struct {
	EndPoint  string
	ModelName string
	Client    *http.Client
}

func NewLLMService(baseUrl string, Client *http.Client) *LLMService {
	return &LLMService{
		EndPoint:  baseUrl + "/api/generate",
		ModelName: "tinyllama",
		Client:    Client,
	}
}

// GenerateResponseWithoutChunks generates a response using only the LLM without any context chunks
func (llm *LLMService) GenerateResponseWithoutChunks(question string) (string, error) {
	prompt := fmt.Sprintf("Respond to this prompt: %s", question)

	generatedResponse, err := llm.generateResponse(prompt)
	return generatedResponse, err
}

// GenerateResponse generates a response using the LLM with provided context chunks
func (llm *LLMService) GenerateResponse(question string, chunks []string) (string, error) {

	data := ""
	for i, chunk := range chunks {
		data += fmt.Sprintf("Chunk %d: %s\n\n", i+1, chunk)
	}

	prompt := fmt.Sprintf("Using this data: %s. Respond to this prompt: %s",
		data,
		question,
	)
	/*
		prompt := fmt.Sprintf("Answer the question using ONLY the information from this context: %s.Respond to this request concisely and directly: %s",
			data,
			question,
		)
	*/
	generatedResponse, err := llm.generateResponse(prompt)
	return generatedResponse, err
}

// generateResponse sends the prompt to the LLM and returns the generated response
func (llm *LLMService) generateResponse(prompt string) (string, error) {

	reqBody := map[string]interface{}{
		"model":  "tinyllama",
		"prompt": prompt,
		"stream": false,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := llm.Client.Post(
		llm.EndPoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //for memory leak

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}
