package models

type QA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type RetrievalEvalData struct {
	Question       string  `json:"question"`
	ExpextedAnswer string  `json:"answer"`
	RelevantChunks []Chunk `json:"relevantChunks"`
}

type RetrievalTestCaseResult struct {
	Question          string
	ExpectedAnswer    string
	ExpectedChunkIDs  []int
	RetrievedChunkIDs []int
	Precision         float64
	Recall            float64
	F1                float64
}

type RetrievalEvaluationResult struct {
	TestCaseResults []RetrievalTestCaseResult
	AvgPrecision    float64
	AvgRecall       float64
	AvgF1           float64
}

type GenerationEvaluationCase struct {
	Question             string
	GroundTruth          string
	GeneratedAnswer      string
	GroundTruthEmbedding []float32
	GeneratedEmbedding   []float32
	SimilarityScore      float64
}

type GenerationEvaluationResult struct {
	TestCaseResults    []GenerationEvaluationCase
	AvgSimilarityScore float64
}
