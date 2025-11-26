package generators

import "net/http"

type GeneratorManager struct {
	EndPoint  string
	ModelName string
	Client    *http.Client
}

type LLMResult struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type OllamaRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options"`
}
