package chunkers

import (
	"log"
	"strings"
)

// NewChunker creates and returns a new ChunkConfig
func NewChunker(chunkSize int, chunkOverlap int) *ChunkManager {
	return &ChunkManager{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

// ChunkText splits the input text into chunks based on the config.Chunk
func (config ChunkManager) ChunkText(text string) []Chunk {
	log.Println("Chunking is started...")
	words := strings.Fields(text)
	log.Printf(" Document length: %d words", len(words))

	var chunks []Chunk
	chunkID := 0

	for i := 0; i < len(words); i += (config.ChunkSize - config.ChunkOverlap) {

		end := i + config.ChunkSize

		//for last chunk
		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[i:end], " ")

		chunks = append(chunks, Chunk{ID: chunkID, Text: chunk})
		chunkID += 1

		if end == len(words) {
			break
		}
	}

	log.Printf(" Chunk size: %d", len(chunks))
	log.Println("Exiting ChunkText")
	return chunks
}
