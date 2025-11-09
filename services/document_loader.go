package services

import (
	"fmt"
	"log"
	"os"
)

// LoadDocument reads the document and returns its content as a string
func loadDocument() (string, error) {

	data, err := os.ReadFile(bookPath)
	if err != nil {
		return "", fmt.Errorf("failed to load document: %w", err)
	}

	text := string(data)
	log.Println("Document was loaded!")

	return text, nil
}
