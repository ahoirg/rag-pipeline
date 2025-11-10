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
