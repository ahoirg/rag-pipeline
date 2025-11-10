package models

type EvalData struct {
	Context string `json:"context"`
	QAS     []QA   `json:"qas"`
}

type QA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type EvalCase struct {
	Question             string
	ExpectedAnswer       string
	GroundTruthContextID int
}

type EvalResult struct {
	EvalCase        EvalCase
	SourceChunks    []string
	GeneratedAnswer string
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
