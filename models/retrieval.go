package models

type RetrievalResult struct {
	ChunkID uint64
	Text    string
	Score   float32 // Cosine similarity score
}
