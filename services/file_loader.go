package services

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
