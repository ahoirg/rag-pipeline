package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rag-pipeline/models"
)

// LoadDocument reads the document and returns its content as a string
func LoadDocument(path string) (string, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("utils.go|LoadDocument: failed to load document: %w", err)
	}

	text := string(data)
	log.Println("utils.go|LoadDocument: Document was loaded!")

	return text, nil
}

// CalculateTruePositive returns the count of retrieved chunk IDs
// that exist in the expected set
func CalculateTruePositive(expected []int, retrieval []int) int {
	var truePositive int = 0

	for _, retrievedChunkID := range retrieval {

		for _, expectedChunkID := range expected {

			if retrievedChunkID == expectedChunkID {
				truePositive++
				break
			}

		}

	}

	return truePositive
}

// Average returns the average of a float64 slice
func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	total := 0.0
	for _, v := range values {
		total += v
	}

	return total / float64(len(values))
}

func StoreChunk(chunk *[]models.Chunk) error {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		return fmt.Errorf("error occurred during marshalling: %w", err)
	}

	file, err := os.Create("chunk.json")
	if err != nil {
		return fmt.Errorf("error occurred during file creating: %s", err.Error())
	}

	if _, err := file.Write(jsonChunk); err != nil {
		return fmt.Errorf("error occurred during jsonChunk writing: %s", err.Error())
	}
	return nil
}
