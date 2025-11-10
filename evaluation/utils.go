package evaluation

import (
	"fmt"
	"log"
	"os"
)

// LoadDocument reads the document and returns its content as a string
func LoadDocument(path string) (string, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load document: %w", err)
	}

	text := string(data)
	log.Println("Document was loaded!")

	return text, nil
}

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
