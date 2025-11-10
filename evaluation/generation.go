package evaluation

import (
	"encoding/json"
	"fmt"
	"rag-pipeline/models"

	"github.com/drewlanenga/govector"
)

func (eval *Evaluator) EvaluateGeneration(generationDataPath string) (*models.GenerationEvaluationResult, error) {
	qas, err := eval.GetGeneratorQA(generationDataPath)
	if err != nil {
		return nil, err
	}

	testCases, err := eval.GetGeneratorResponse(qas)
	if err != nil {
		return nil, err
	}

	err = eval.GetEmbeddedResults(&testCases)
	if err != nil {
		return nil, err
	}

	err = eval.CalculateSimilarityScores(&testCases)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	for _, tc := range testCases {
		totalScore += tc.SimilarityScore
	}

	result := &models.GenerationEvaluationResult{
		TestCaseResults:    testCases,
		AvgSimilarityScore: totalScore / float64(len(testCases)),
	}

	return result, nil
}

func (eval *Evaluator) GetGeneratorQA(generationDataPath string) ([]models.QA, error) {

	text, err := LoadDocument(generationDataPath)
	if err != nil {
		return nil, fmt.Errorf("generation.go |failed to load QA JSON: %w", err)
	}

	var qaData []models.QA
	if err := json.Unmarshal([]byte(text), &qaData); err != nil {
		return nil, fmt.Errorf("generation.go |failed to parse QA JSON: %w", err)
	}

	return qaData, nil
}

func (eval *Evaluator) GetGeneratorResponse(qaData []models.QA) ([]models.GenerationEvaluationCase, error) {
	var evaluationCase []models.GenerationEvaluationCase

	for _, qa := range qaData {
		generatedAnswer, _, err := eval.RAGService.GenerateResponse(qa.Question)
		if err != nil {
			return nil, fmt.Errorf("generation.go |failed to generate response: %w", err)
		}

		evaluationCase = append(evaluationCase, models.GenerationEvaluationCase{
			Question:        qa.Question,
			GroundTruth:     qa.Answer,
			GeneratedAnswer: generatedAnswer,
		})
	}

	return evaluationCase, nil
}

func (eval *Evaluator) GetEmbeddedResults(results *[]models.GenerationEvaluationCase) error {
	var groundTruths []string
	var generatedAnswers []string

	for _, r := range *results {
		groundTruths = append(groundTruths, r.GroundTruth)
		generatedAnswers = append(generatedAnswers, r.GeneratedAnswer)
	}

	groundTruthEmbeddings, err := eval.RAGService.Embedder.EmbedChunks(groundTruths)
	if err != nil {
		return fmt.Errorf("generation.go|failed to embed ground truths: %w", err)
	}

	generatedEmbeddings, err := eval.RAGService.Embedder.EmbedChunks(generatedAnswers)
	if err != nil {
		return fmt.Errorf("generation.go|failed to embed generated answers: %w", err)
	}

	for i := range *results {
		(*results)[i].GroundTruthEmbedding = groundTruthEmbeddings[i]
		(*results)[i].GeneratedEmbedding = generatedEmbeddings[i]
	}

	return nil
}

func (eval *Evaluator) CalculateSimilarityScores(results *[]models.GenerationEvaluationCase) error {
	for i := range *results {

		// TODO
		// maybe we can store embedding in float64 ?
		// need to analyze
		// writemanuel ???
		vec1 := make([]float64, len((*results)[i].GroundTruthEmbedding))
		vec2 := make([]float64, len((*results)[i].GeneratedEmbedding))

		for j := range vec1 {
			vec1[j] = float64((*results)[i].GroundTruthEmbedding[j])
			vec2[j] = float64((*results)[i].GeneratedEmbedding[j])
		}

		var err error
		(*results)[i].SimilarityScore, err = govector.Cosine(vec1, vec2)
		if err != nil {
			return err
		}

	}
	return nil
}
