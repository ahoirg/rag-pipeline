package evaluation

import (
	"encoding/json"
	"fmt"
	"log"
	"rag-pipeline/models"
	"rag-pipeline/services"
)

const (
	evaluationDataPath       = "evalData/notre_dame_eval.json"
	evaluationSourceDataPath = "evalData/notre_dame_eval.json"
	evalCollectionName       = "eval_collection"
)

type Evaluator struct {
	RAGService   *services.RAGService
	EvalCases    []models.EvalCase
	IsIntialized bool
}

// NewEvaluator initializes the Evaluator with its own RAGService
func NewEvaluator() *Evaluator {
	return &Evaluator{
		RAGService:   services.NewRAGService(evalCollectionName),
		IsIntialized: false,
	}
}

func (eval *Evaluator) GetEvaluationRessults() (string, error) {
	if err := eval.initializeEval(); err != nil {
		log.Println("could not create Evaluation collection.", err)
		return "", err
	}
	return "", nil
}

func (eval *Evaluator) Evaluation() {
	if err := eval.initializeEval(); err != nil {
		log.Println("could not create Evaluation collection.", err)
	}
}

func (eval *Evaluator) initializeEval() error {
	if eval.IsIntialized {
		return nil
	}

	text, err := services.LoadDocument(evaluationSourceDataPath)
	if err != nil {
		return fmt.Errorf("failed prepareQdrantDB: %w", err)
	}

	if err := eval.RAGService.StoreData(text); err != nil {
		return fmt.Errorf("failed prepareQdrantDB: %w", err)
	}

	evalData, err := loadEvaluationData(evaluationDataPath)
	if err != nil {
		return err
	}

	eval.EvalCases = generateTestCases(evalData)

	eval.initializeEval()
	return nil
}

func loadEvaluationData(filePath string) ([]models.EvalData, error) {

	jsonString, err := services.LoadDocument(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load evaluation file: %w", err)
	}

	var evalData []models.EvalData
	if err := json.Unmarshal([]byte(jsonString), &evalData); err != nil {
		return nil, fmt.Errorf("failed to parse evaluation JSON: %w", err)
	}

	log.Printf("Loaded %d contexts with questions\n", len(evalData))

	return evalData, nil
}

func generateTestCases(evalData []models.EvalData) []models.EvalCase {
	var testCases []models.EvalCase

	for contextID, data := range evalData {
		for _, qa := range data.QAS {
			testCases = append(testCases, models.EvalCase{
				Question:             qa.Question,
				ExpectedAnswer:       qa.Answer,
				GroundTruthContextID: contextID,
			})
		}
	}

	log.Printf("Generated %d test cases\n", len(testCases))
	return testCases
}
