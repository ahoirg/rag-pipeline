package models

type RetrievalResult struct {
	ChunkID int
	Text    string
	Score   float32 // Cosine similarity score
}
