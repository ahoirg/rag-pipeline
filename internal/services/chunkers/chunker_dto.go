package chunkers

type ChunkManager struct {
	ChunkSize    int
	ChunkOverlap int
}

type Chunk struct {
	ID   int    `json:"chunkID"`
	Text string `json:"text"`
}
