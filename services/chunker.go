package services

import (
	"log"
	"strings"
)

type ChunkConfig struct {
	ChunkSize    int
	ChunkOverlap int
}

func NewChunker(chunkSize int, chunkOverlap int) *ChunkConfig {
	return &ChunkConfig{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

// TODO : basicly splitting into chunks for now
// ChunkText splits the input text into chunks based on the config(ChunkConfig)
func (config ChunkConfig) ChunkText(text string) []string {

	words := strings.Fields(text)
	log.Printf(" Document length: %d words", len(words))

	chunks := []string{}

	// Iterate through the words, creating chunks
	// yhe step size is calculated by ChunkSize - ChunkOverlap
	for i := 0; i < len(words); i += (config.ChunkSize - config.ChunkOverlap) {
		end := i + config.ChunkSize
		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)

		if end == len(words) {
			break
		}
	}

	log.Printf(" Chunk size: %d", len(chunks))
	return chunks
}
