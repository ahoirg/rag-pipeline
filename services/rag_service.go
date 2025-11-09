package services

import (
	"fmt"
	"log"
	"net/http"
	"rag-pipeline/db"
	"rag-pipeline/models"
	"slices"
	"time"
)

type RAGService struct {
	Chunker   *ChunkConfig
	Embedder  *OllamaEmbedder
	QdrantDB  *db.QdrantDatabase
	Generator *LLMService
}

// TODO: move to config file
const (
	bookPath               = "data/treasure_island.txt"
	chunkSize              = 300
	chunkOverlap           = 30
	qdrantHost             = "localhost"
	qdrantPort             = 6334
	collectionName         = "treasure_island"
	embedindModelDimention = 768 // Nomic embedder dimension
	embeddingModel         = "nomic-embed-text"
	llomaBaseURL           = "http://localhost:11434"
)

// NewRAGService initializes the RAG service by setting up the Qdrant client and preparing the vector database
func NewRAGService() (*RAGService, error) {

	log.Println(" Starting RAG Service...")

	ragService := &RAGService{
		Chunker:  NewChunker(chunkSize, chunkOverlap),
		Embedder: NewOllamaEmbedder(llomaBaseURL, embeddingModel),
		QdrantDB: db.NewQdrantDatabase(qdrantHost, qdrantPort, collectionName),
		Generator: NewLLMService(llomaBaseURL, &http.Client{
			Timeout: 5 * time.Minute,
		}),
	}
	if err := ragService.prepareVectorDatabase(); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	log.Println(" RAG Service is working.")
	return ragService, nil
}

func (r *RAGService) GenerateResponse(question string) (string, error) {
	retrievalResult, err := r.RetrieveRelevantChunks(question, 3)
	if err != nil {
		return "", err
	}

	var chunks []string
	for _, res := range retrievalResult {
		chunks = append(chunks, res.Text)
	}

	generatedResponse, err := r.Generator.GenerateResponse(question, chunks)

	return generatedResponse, err
}

func (r *RAGService) GenerateResponseWithoutChunks(question string) (string, error) {
	generatedResponse, err := r.Generator.GenerateResponseWithoutChunks(question)
	return generatedResponse, err
}

func (r *RAGService) RetrieveRelevantChunks(query string, topK int) ([]models.RetrievalResult, error) {

	queryEmbedding, err := r.Embedder.EmbedQuery(query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	searchResult, err := r.QdrantDB.QueryQdrant(queryEmbedding, uint64(topK))
	if err != nil {
		return nil, fmt.Errorf("failed to query Qdrant: %w", err)
	}

	var results []models.RetrievalResult
	for _, point := range searchResult {
		results = append(results, models.RetrievalResult{
			ChunkID: point.Id.GetNum(),
			Text:    point.Payload["text"].GetStringValue(),
			Score:   point.Score,
		})
	}

	return results, nil
}

// prepareVectorDatabase checks if the required collection exists in Qdrant
// if not, it creates collection and insert chunked and embedded data
func (r *RAGService) prepareVectorDatabase() error {

	collections, err := r.QdrantDB.GetQdrantCollectionNames()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	if len(collections) > 0 && slices.Contains(collections, "treasure_island") {
		return nil
	}

	chunks, embeddings, err := r.chunkandEmbed()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	err = r.QdrantDB.CreateQdrantCollection(uint64(embedindModelDimention))
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	err = r.QdrantDB.AddVectorsToQdrant(chunks, embeddings)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	return nil
}

// chunkandEmbed loads the document, chunks it and generates embeddings for the chunks
func (r *RAGService) chunkandEmbed() ([]string, [][]float32, error) {

	text, err := loadDocument()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, nil, err
	}

	chunks := r.Chunker.ChunkText(text)
	if len(chunks) == 0 {
		return nil, nil, fmt.Errorf("chunking failed: no chunks were created from the given text")
	}

	embeddings, err := r.Embedder.EmbedChunks(chunks)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, nil, err
	}

	return chunks, embeddings, nil
}
