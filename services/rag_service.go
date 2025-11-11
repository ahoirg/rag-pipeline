package services

import (
	"fmt"
	"net/http"
	"rag-pipeline/db"
	"rag-pipeline/models"
	"time"
)

type RAGService struct {
	Chunker   *ChunkConfig
	Embedder  *OllamaEmbedder
	QdrantDB  *db.QdrantDatabase
	Generator *LLMService
	Config    *models.Config
}

// NewRAGService initializes the RAG service by setting up the Qdrant client and preparing the vector database
func NewRAGService(config *models.Config, collectionName string) *RAGService {
	return &RAGService{
		Chunker:  NewChunker(config.Chunk.Size, config.Chunk.Overlap),
		Embedder: NewOllamaEmbedder(config.Ollama.BaseURL, config.Embedding.ModelName, config.Embedding.Endpoint),
		Generator: NewLLMService(config.Ollama.BaseURL, config.Generator.Endpoint, config.Generator.ModelName, &http.Client{
			Timeout: 5 * time.Minute,
		}),
		QdrantDB: db.NewQdrantDatabase(config.Qdrant.Host, config.Qdrant.Port, collectionName),
		Config:   config,
	}
}

// StoreData sends the given text data to the vector database
func (r *RAGService) StoreData(text string) error {
	return r.prepareVectorDatabase(text)
}

// GenerateResponse retrieves the most relevant chunks for the given question,
// sends them with the query to the generator model and returns the generated answer
func (r *RAGService) GenerateResponse(question string) (string, []string, error) {
	retrievalResult, err := r.RetrieveRelevantChunks(question, 3)
	if err != nil {
		return "", nil, err
	}

	var chunks []string
	for _, res := range retrievalResult {
		chunks = append(chunks, res.Text)
	}

	generatedResponse, err := r.Generator.GenerateResponse(question, chunks)

	return generatedResponse, chunks, err
}

// GenerateResponseWithoutChunks sends the given question directly to the the generator model
// it returns the generated answer
func (r *RAGService) GenerateResponseWithoutChunks(question string) (string, error) {
	return r.Generator.GenerateResponseWithoutChunks(question)
}

// RetrieveRelevantChunks retrieves the most relevant chunks for the given query
func (r *RAGService) RetrieveRelevantChunks(query string, topK int) ([]models.RetrievalResult, error) {

	queryEmbedding, err := r.Embedder.EmbedQuery(query)
	if err != nil {
		return nil, fmt.Errorf("RetrieveRelevantChunks: failed to embed query: %w", err)
	}

	searchResult, err := r.QdrantDB.QueryQdrant(queryEmbedding, uint64(topK))
	if err != nil {
		return nil, fmt.Errorf("RetrieveRelevantChunks: failed to query Qdrant: %w", err)
	}

	var results []models.RetrievalResult
	for _, point := range searchResult {
		results = append(results, models.RetrievalResult{
			ChunkID: int(point.Payload["id"].GetIntegerValue()),
			Text:    point.Payload["text"].GetStringValue(),
			Score:   point.Score,
		})
	}

	return results, nil
}

// prepareVectorDatabase insert chunked and embedded data to db
func (r *RAGService) prepareVectorDatabase(text string) error {

	//Chunks
	chunks := r.Chunker.ChunkText(text)
	if len(chunks) == 0 {
		return fmt.Errorf("rag_serivece| prepareVectorDatabase: chunking failed: no chunks were created from the given text")
	}

	//prepare chunks for embeddings
	chunk_texts := make([]string, len(chunks)) // 'make' for fast, direct indext assignment and no allocation
	for i, chunk := range chunks {
		chunk_texts[i] = chunk.Text
	}

	//embedding
	embeddings, err := r.Embedder.EmbedChunks(chunk_texts)
	if err != nil {
		return fmt.Errorf("rag_serivece| prepareVectorDatabase: %w", err)
	}

	//TODO  handle it in db depending on qdrantdb client
	err = r.QdrantDB.CreateQdrantCollection(uint64(r.Config.Embedding.ModelDimension))
	if err != nil {
		return fmt.Errorf("rag_serivece| prepareVectorDatabase: %w", err)
	}

	//stores vectors in db
	err = r.QdrantDB.AddVectorsToQdrant(chunks, embeddings)
	if err != nil {
		return fmt.Errorf("rag_serivece| prepareVectorDatabase: %w", err)
	}

	return nil
}
