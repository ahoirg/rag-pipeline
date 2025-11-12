package evaluation

import (
	"encoding/json"
	"fmt"
	"rag-pipeline/models"
	"rag-pipeline/utils"

	"github.com/drewlanenga/govector"
)

// EvaluateGeneration runs the full evaluation pipeline for generated responses
func (eval *Evaluator) EvaluateGeneration(generationDataPath string) (*models.GenerationEvaluationResult, error) {
	qas, err := eval.loadGeneratorQA(generationDataPath)
	if err != nil {
		return nil, err
	}

	generationCases, err := eval.getGeneratorResponse(qas)
	if err != nil {
		return nil, err
	}

	err = eval.getEmbeddedResults(&generationCases)
	if err != nil {
		return nil, err
	}

	err = eval.calculateSimilarityScores(&generationCases)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	for _, tc := range generationCases {
		totalScore += tc.SimilarityScore
	}

	result := &models.GenerationEvaluationResult{
		TestCaseResults:    generationCases,
		AvgSimilarityScore: totalScore / float64(len(generationCases)),
	}

	return result, nil
}

// getGeneratorResponse generates answers for each QA pair and returns them with generated answers
func (eval *Evaluator) getGeneratorResponse(qaData []models.QA) ([]models.GenerationEvaluationCase, error) {
	var evaluationCase []models.GenerationEvaluationCase

	for _, qa := range qaData {
		generatedAnswer, chunks, err := eval.RAGService.GenerateResponse(qa.Question)
		if err != nil {
			return nil, fmt.Errorf("generation.go |failed to generate response: %w", err)
		}

		evaluationCase = append(evaluationCase, models.GenerationEvaluationCase{
			Question:        qa.Question,
			GroundTruth:     qa.Answer,
			GeneratedAnswer: generatedAnswer,
			SourceChunks:    chunks,
		})
	}

	return evaluationCase, nil
}

// getEmbeddedResults generates embeddings for ground truths and generated answers
func (eval *Evaluator) getEmbeddedResults(results *[]models.GenerationEvaluationCase) error {
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

// calculateSimilarityScores computes cosine similarity between ground truth and generated embeddings
func (eval *Evaluator) calculateSimilarityScores(results *[]models.GenerationEvaluationCase) error {
	for i := range *results {

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

// loadGeneratorQA reads QA data from a JSON file and returns it as a slice of models.QA
func (eval *Evaluator) loadGeneratorQA(generationDataPath string) ([]models.QA, error) {

	text, err := utils.LoadDocument(generationDataPath)
	if err != nil {
		return nil, fmt.Errorf("generation.go |failed to load QA JSON: %w", err)
	}

	var qaData []models.QA
	if err := json.Unmarshal([]byte(text), &qaData); err != nil {
		return nil, fmt.Errorf("generation.go |failed to parse QA JSON: %w", err)
	}

	return qaData, nil
}
