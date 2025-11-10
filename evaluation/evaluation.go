package evaluation

import (
	"fmt"
	"log"
	"rag-pipeline/models"
	"rag-pipeline/services"
)

// TODO: moveto config
const (
	evaluation_retrievalDataPath = "evalData/retrieval/notre_dame_qa_chunks.json"
	evaluationSourceDataPath     = "evalData/notre_dame_contexts.txt"
	evalCollectionName           = "eval_collection"
)

type Evaluator struct {
	RAGService                *services.RAGService
	RetrievalEvaluationResult *models.RetrievalEvaluationResult
	IsIntialized              bool
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		RAGService:                services.NewRAGService(evalCollectionName),
		RetrievalEvaluationResult: nil,   //
		IsIntialized:              false, // insert evaluation data into the database only once
	}
}

func (eval *Evaluator) GetRetrievalEvaluateResult() (*models.RetrievalEvaluationResult, error) {
	if eval.RetrievalEvaluationResult != nil {
		return eval.RetrievalEvaluationResult, nil
	}

	if err := eval.prepareEvalData(); err != nil {
		log.Println("evaluation.go| could not create evaluation collection.", err)
		return nil, err
	}

	retrievalEvaluationResult, err := eval.EvaluateRetrieval(evaluation_retrievalDataPath)
	if err != nil {
		log.Println("evaluation.go| could not run retrieval evaluation.", err)
		return nil, err
	}

	eval.RetrievalEvaluationResult = retrievalEvaluationResult
	return eval.RetrievalEvaluationResult, nil
}

func (eval *Evaluator) prepareEvalData() error {
	if eval.IsIntialized {
		return nil
	}

	text, err := LoadDocument(evaluationSourceDataPath)
	if err != nil {
		return fmt.Errorf("evaluation.go|failed prepareQdrantDB: %w", err)
	}

	if err := eval.RAGService.StoreData(text); err != nil {
		return fmt.Errorf("evaluation.go|failed prepareQdrantDB: %w", err)
	}

	eval.IsIntialized = true
	return nil
}
