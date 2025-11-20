package chunkers

import (
	"encoding/json"
	"fmt"
	"os"
)

// StoreChunk stores chunks in the root directory as a json
func StoreChunk(chunk *[]Chunk) error {
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
