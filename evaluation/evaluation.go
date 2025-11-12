package evaluation

import (
	"fmt"
	"log"
	"rag-pipeline/models"
	"rag-pipeline/services"
	"rag-pipeline/utils"
)

type Evaluator struct {
	RAGService                 *services.RAGService
	RetrievalEvaluationResult  *models.RetrievalEvaluationResult
	GenerationEvaluationResult *models.GenerationEvaluationResult
	Config                     *models.Config
}

// NewChunker creates and returns a new Evaluator
func NewEvaluator(config *models.Config) (*Evaluator, error) {
	ragService, err := services.NewRAGService(config, config.Evaluation.CollectionName)
	if err != nil {
		return nil, fmt.Errorf("evaluation.go| NewEvaluator: initialization error %w", err)
	}

	return &Evaluator{
		RAGService: ragService,

		// Stores results to avoid recalculating on the same data with
		// the same parameters once the evaluation has been performed
		RetrievalEvaluationResult:  nil,
		GenerationEvaluationResult: nil,

		Config: config,
	}, nil
}

// GetRetrievalEvaluateResult returns the retrieval evaluation result
func (eval *Evaluator) GetRetrievalEvaluateResult() (*models.RetrievalEvaluationResult, error) {
	if eval.RetrievalEvaluationResult != nil {
		return eval.RetrievalEvaluationResult, nil
	}

	if err := eval.prepareEvalData(); err != nil {
		log.Println("evaluation.go| could not create evaluation collection.", err)
		return nil, err
	}

	retrievalEvaluationResult, err := eval.EvaluateRetrieval(eval.Config.Evaluation.RetrievalDataPath)
	if err != nil {
		log.Println("evaluation.go| could not run retrieval evaluation.", err)
		return nil, err
	}

	eval.RetrievalEvaluationResult = retrievalEvaluationResult
	return eval.RetrievalEvaluationResult, nil
}

// GetGenerationEvaluateResult returns the generation evaluation result
func (eval *Evaluator) GetGenerationEvaluateResult() (*models.GenerationEvaluationResult, error) {
	if eval.GenerationEvaluationResult != nil {
		return eval.GenerationEvaluationResult, nil
	}

	if err := eval.prepareEvalData(); err != nil {
		log.Println("evaluation.go| could not create evaluation collection.", err)
		return nil, err
	}

	generationEvaluationResult, err := eval.EvaluateGeneration(eval.Config.Evaluation.GenerationDataPath)
	if err != nil {
		log.Println("evaluation.go| could not run retrieval evaluation.", err)
		return nil, err
	}

	eval.GenerationEvaluationResult = generationEvaluationResult
	return eval.GenerationEvaluationResult, nil
}

// prepareEvalData loads the evaluation source data and stores it in the vector database
func (eval *Evaluator) prepareEvalData() error {

	// If either of the evaluations has already been computed,
	// it means the test data has been loaded previously.
	if eval.RetrievalEvaluationResult != nil || eval.GenerationEvaluationResult != nil {
		return nil
	}

	text, err := utils.LoadDocument(eval.Config.Evaluation.SourceDataPath)
	if err != nil {
		return fmt.Errorf("evaluation.go|failed prepareQdrantDB: %w", err)
	}

	if err := eval.RAGService.StoreData(text); err != nil {
		return fmt.Errorf("evaluation.go|failed prepareQdrantDB: %w", err)
	}

	return nil
}
