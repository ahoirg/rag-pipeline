// retrieval.go
package evaluation

import (
	"encoding/json"
	"fmt"
	"rag-pipeline/models"
	"rag-pipeline/utils"
)

// EvaluateRetrieval runs the full evaluation pipeline for retrieval evaluation
func (eval *Evaluator) EvaluateRetrieval(retrievalDataPath string) (*models.RetrievalEvaluationResult, error) {
	evalData, err := loadQuestions(retrievalDataPath)
	if err != nil {
		return nil, err
	}

	var testCaseResults []models.RetrievalTestCaseResult
	for _, data := range evalData {
		retrievedChunks, err := eval.RAGService.RetrieveRelevantChunks(data.Question, eval.Config.Retrieval.TopK)
		if err != nil {
			return nil, err
		}

		retrievedChunkIDs := make([]int, len(retrievedChunks))
		for i, chunk := range retrievedChunks {
			retrievedChunkIDs[i] = chunk.ChunkID
		}

		relevantIDs := make([]int, len(data.RelevantChunks))
		for i, chunk := range data.RelevantChunks {
			relevantIDs[i] = chunk.ID
		}

		testCaseResults = append(testCaseResults, models.RetrievalTestCaseResult{
			Question:          data.Question,
			ExpectedAnswer:    data.ExpextedAnswer,
			ExpectedChunkIDs:  relevantIDs,
			RetrievedChunkIDs: retrievedChunkIDs,
		})
	}

	return calculateRetrievalMetricResults(testCaseResults), nil
}

// calculateRetrievalMetricResults computes precision, recall, and F1 scores for each retrieval test case
// and return them as a models.RetrievalEvaluationResult
func calculateRetrievalMetricResults(testCaseResults []models.RetrievalTestCaseResult) *models.RetrievalEvaluationResult {
	var precisions []float64
	var recalls []float64
	var f1s []float64

	duplicatedTestCaseResults := make([]models.RetrievalTestCaseResult, len(testCaseResults))
	copy(duplicatedTestCaseResults, testCaseResults)

	for i, tc := range testCaseResults {
		truePositives := utils.CalculateTruePositive(tc.ExpectedChunkIDs, tc.RetrievedChunkIDs)
		falsePositive := len(tc.RetrievedChunkIDs) - truePositives
		falseNegative := len(tc.ExpectedChunkIDs) - truePositives

		precision := float64(truePositives) / float64(truePositives+falsePositive)
		recall := float64(truePositives) / float64(truePositives+falseNegative)

		var f1 float64
		if precision+recall > 0 {
			f1 = (2 * precision * recall) / (precision + recall)
		} else {
			f1 = 0
		}

		duplicatedTestCaseResults[i].Precision = precision
		duplicatedTestCaseResults[i].Recall = recall
		duplicatedTestCaseResults[i].F1 = f1

		precisions = append(precisions, precision)
		recalls = append(recalls, recall)
		f1s = append(f1s, f1)
	}

	return &models.RetrievalEvaluationResult{
		TestCaseResults: duplicatedTestCaseResults,
		AvgPrecision:    utils.Average(precisions),
		AvgRecall:       utils.Average(recalls),
		AvgF1:           utils.Average(f1s),
	}
}

// loadQuestions reads QA data from a JSON file and returns it as a slice of models.RetrievalEvalData
func loadQuestions(filePath string) ([]models.RetrievalEvalData, error) {
	jsonString, err := utils.LoadDocument(filePath)
	if err != nil {
		return nil, fmt.Errorf("retrieval.go| failed to LoadQuestions file: %w", err)
	}

	var retrievalEvalData []models.RetrievalEvalData
	if err := json.Unmarshal([]byte(jsonString), &retrievalEvalData); err != nil {
		return nil, fmt.Errorf("retrieval.go| failed to LoadQuestions JSON: %w", err)
	}
	println("retrieval.go| Questions were loaded!")
	return retrievalEvalData, nil
}
