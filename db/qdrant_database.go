package db

import (
	"context"
	"fmt"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

// NewQdrantClient creates and returns a new Qdrant client
func NewQdrantClient(qdrantHost string, qdrantPort int) (*qdrant.Client, error) {
	return qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantPort,
	})
}

// GetQdrantCollectionNames gets the names of all collections in the Qdrant database
func GetQdrantCollectionNames(client *qdrant.Client) ([]string, error) {
	names, err := client.ListCollections(context.Background())
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return names, nil
}

// CreateQdrantCollection creates a new collection in Qdrant with the collectionName and vector size
func CreateQdrantCollection(client *qdrant.Client, collectionName string, vectorSize uint64) error {
	return client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

// AddVectorsToQdrant adds the given chunks and their corresponding embeddings to the collection
func AddVectorsToQdrant(client *qdrant.Client, collectionName string, chunks []string, embeddings [][]float32) error {

	var points []*qdrant.PointStruct

	for i := 0; i < len(chunks); i++ {
		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i)),
			Vectors: qdrant.NewVectors(embeddings[i]...),
			Payload: qdrant.NewValueMap(map[string]any{
				"text": chunks[i],
			}),
		})
	}

	operationInfo, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	if err != nil {
		return err
	}

	fmt.Println(operationInfo)

	return nil
}
