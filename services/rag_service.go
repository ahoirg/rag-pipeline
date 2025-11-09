package services

import (
	"fmt"
	"log"
	"net/http"
	"rag-pipeline/db"
	"rag-pipeline/models"
	"slices"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

type RAGService struct {
	QdrantClient *qdrant.Client
}

// TODO: move to config file
const (
	bookPath               = "data/treasure_island.txt"
	chunkSize              = 300
	chunkOverlap           = 30
	QdrantHost             = "localhost"
	QdrantPort             = 6334
	collectionName         = "treasure_island"
	EmbedindModelDimention = 768 // Nomic embedder dimension
)

// NewRAGService initializes the RAG service by setting up the Qdrant client and preparing the vector database
func NewRAGService() (*RAGService, error) {
	log.Println(" Starting RAG Service...")

	q_client, err := db.NewQdrantClient(QdrantHost, QdrantPort)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	if err := prepareVectorDatabase(q_client); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	ragService := &RAGService{
		QdrantClient: q_client,
	}

	return ragService, nil
}

// prepareVectorDatabase checks if the required collection exists in Qdrant
// if not, it creates collection and insert chunked and embedded data
func prepareVectorDatabase(q_client *qdrant.Client) error {

	collections, err := db.GetQdrantCollectionNames(q_client)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	if len(collections) > 0 && slices.Contains(collections, "treasure_island") {
		return nil
	}

	chunks, embeddings, err := chunkandEmbed()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	err = db.CreateQdrantCollection(q_client, collectionName, uint64(EmbedindModelDimention))
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	err = db.AddVectorsToQdrant(q_client, collectionName, chunks, embeddings)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	return nil
}

// chunkandEmbed loads the document, chunks it and generates embeddings for the chunks
func chunkandEmbed() ([]string, [][]float32, error) {

	text, err := loadDocument()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, nil, err
	}

	// Define chunking configuration
	config := models.ChunkConfig{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}

	chunks := ChunkText(text, config)
	if len(chunks) == 0 {
		return nil, nil, fmt.Errorf("chunking failed: no chunks were created from the given text")
	}

	embedder := models.OllamaEmbedderConfig{
		BaseURL: "http://localhost:11434",
		Model:   "nomic-embed-text",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	embeddings, err := EmbedChunks(chunks, embedder)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, nil, err
	}

	return chunks, embeddings, nil
}
