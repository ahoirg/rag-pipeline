package services

import (
	"log"
	"rag-pipeline/models"
	"strings"
)

type ChunkConfig struct {
	ChunkSize    int
	ChunkOverlap int
}

// NewChunker creates and returns a new ChunkConfig
func NewChunker(chunkSize int, chunkOverlap int) *ChunkConfig {
	return &ChunkConfig{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

// TODO : basicly splitting into chunks for now
// ChunkText splits the input text into chunks based on the config.Chunk
func (config ChunkConfig) ChunkText(text string) []models.Chunk {

	words := strings.Fields(text)
	log.Printf(" Document length: %d words", len(words))

	var chunks []models.Chunk
	chunkID := 0

	// Iterate through the words, creating chunks
	// yhe step size is calculated by ChunkSize - ChunkOverlap
	for i := 0; i < len(words); i += (config.ChunkSize - config.ChunkOverlap) {

		end := i + config.ChunkSize

		//for last chunk
		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[i:end], " ")

		chunks = append(chunks, models.Chunk{ID: chunkID, Text: chunk})
		chunkID += 1

		if end == len(words) {
			break
		}
	}

	log.Printf(" Chunk size: %d", len(chunks))
	return chunks
}
