package db

import (
	"context"
	"fmt"
	"log"
	"rag-pipeline/models"

	"github.com/qdrant/go-client/qdrant"
)

// QdrantDatabase represents a Qdrant database connection
type QdrantDatabase struct {
	Client         *qdrant.Client
	CollectionName string
}

// NewQdrantDatabase creates and returns a new QdrantDatabase instance
func NewQdrantDatabase(qdrantHost string, qdrantPort int, collectionName string) *QdrantDatabase {
	client, err := newQdrantClient(qdrantHost, qdrantPort)
	if err != nil {
		log.Fatalf("Failed to create Qdrant client: %v", err)
	}

	return &QdrantDatabase{
		Client:         client,
		CollectionName: collectionName,
	}
}

// GetQdrantCollectionNames gets the names of all collections in the Qdrant database
func (qdb *QdrantDatabase) GetQdrantCollectionNames() ([]string, error) {
	return qdb.Client.ListCollections(context.Background())
}

func (qdb *QdrantDatabase) CollectionExists() (bool, error) {
	return qdb.Client.CollectionExists(context.Background(), qdb.CollectionName)
}

// CreateQdrantCollection creates a new collection in Qdrant with the collectionName and vector size
func (qdb *QdrantDatabase) CreateQdrantCollection(vectorSize uint64) error {
	return qdb.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: qdb.CollectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

// AddVectorsToQdrant adds the given chunks and their corresponding embeddings to the collection
func (qdb *QdrantDatabase) AddVectorsToQdrant(chunks []models.Chunk, embeddings [][]float32) error {

	var points []*qdrant.PointStruct

	for i := 0; i < len(chunks); i++ {
		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i)),
			Vectors: qdrant.NewVectors(embeddings[i]...),
			Payload: qdrant.NewValueMap(map[string]any{
				"id":   chunks[i].ID,
				"text": chunks[i].Text,
			}),
		})
	}

	operationInfo, err := qdb.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: qdb.CollectionName,
		Points:         points,
	})
	if err != nil {
		return err
	}

	fmt.Println(operationInfo)

	return nil
}

func (qdb *QdrantDatabase) QueryQdrant(queryEmbedding []float32, limit uint64) ([]*qdrant.ScoredPoint, error) {

	searchResult, err := qdb.Client.Query(
		context.Background(),
		&qdrant.QueryPoints{
			CollectionName: qdb.CollectionName,
			Query:          qdrant.NewQuery(queryEmbedding...),
			Limit:          &limit,
			WithPayload:    qdrant.NewWithPayload(true),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("qdrant_database: failed to query Qdrant: %w", err)
	}

	return searchResult, nil
}

func (qdb *QdrantDatabase) DeleteCollection() error {
	return qdb.Client.DeleteCollection(context.Background(), qdb.CollectionName)
}

// newQdrantClient creates and returns a new Qdrant client
func newQdrantClient(qdrantHost string, qdrantPort int) (*qdrant.Client, error) {
	return qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantPort,
	})
}
