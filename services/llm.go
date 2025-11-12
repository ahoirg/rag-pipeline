package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rag-pipeline/models"
	"time"
)

type LLMService struct {
	EndPoint  string
	ModelName string
	Client    *http.Client
}

// NewLLMService creates and returns a new LLMService
func NewLLMService(baseUrl string, endpoint string, modelName string) *LLMService {
	return &LLMService{
		EndPoint:  baseUrl + endpoint,
		ModelName: modelName,
		Client: &http.Client{
			Timeout: 120 * time.Second},
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

	prompt := fmt.Sprintf(`We have provided context information below.
---------------------
%s
---------------------
Answer the question with only the essential information. Just write the answer to the Question
Question: %s
Answer: \
`,
		data,
		question,
	)

	generatedResponse, err := llm.generateResponse(prompt)
	return generatedResponse, err
}

// generateResponse sends the prompt to the LLM and returns the generated response
func (llm *LLMService) generateResponse(prompt string) (string, error) {

	reqBody := models.OllamaRequest{
		Model:  llm.ModelName,
		Prompt: prompt,
		Stream: false,
		Options: map[string]any{
			"temperature": 0.1,
		},
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

	var result models.LLMResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}
