package models

type EmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbedQueryRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embedding  []float32   `json:"embedding"`
	Embeddings [][]float32 `json:"embeddings"`
}
