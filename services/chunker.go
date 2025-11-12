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

// ChunkText splits the input text into chunks based on the config.Chunk
func (config ChunkConfig) ChunkText(text string) []models.Chunk {
	log.Println("Chunking is started...")
	words := strings.Fields(text)
	log.Printf(" Document length: %d words", len(words))

	var chunks []models.Chunk
	chunkID := 0

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

	//utils.StoreChunk(&chunks)

	log.Printf(" Chunk size: %d", len(chunks))
	log.Println("Exiting ChunkText")
	return chunks
}

// I do not delete them, maybe we will use them again
/*
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
*/
/*
	segmenter := sentencizer.NewSegmenter("en")
	sentences := segmenter.Segment(text)
	log.Println(len(sentences))

	for i := 0; i < len(sentences); i += (config.ChunkSize - config.ChunkOverlap) {
		end := i + config.ChunkSize

		if end > len(sentences) {
			end = len(sentences)
		}

		chunk := strings.Join(sentences[i:end], " ")
		chunks = append(chunks, models.Chunk{ID: chunkID, Text: chunk})
		chunkID += 1

		if end == len(sentences) {
			break
		}
	}
*/ /*
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatalf("Error creating document: %v", err)
	}

	sentences := make([]string, len(doc.Sentences()))
	for i, s := range doc.Sentences() {
		sentences[i] = s.Text
	}

	var chunks []models.Chunk
	chunkID := 0
	for i := 0; i < len(sentences); i += (config.ChunkSize - config.ChunkOverlap) {
		end := i + config.ChunkSize

		if end > len(sentences) {
			end = len(sentences)
		}

		chunk := strings.Join(sentences[i:end], " ")
		chunks = append(chunks, models.Chunk{ID: chunkID, Text: chunk})
		chunkID += 1

		if end == len(sentences) {
			break
		}
	}*/
