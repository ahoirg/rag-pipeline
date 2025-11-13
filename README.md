[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# Rag Pipeline
It is a modular Retrieval‑Augmented Generation (RAG) pipeline written in Go. It ingests documents, chunks text, generates embeddings, stores vectors and answers questions via retrieval + generation.

## Table of Contents 📑
- [1. Core Features](#1-core-features)
- [2. Tech Stack ](#2-tech-stack)
- [3.1 Docker Container Usage Instructions](#31-docker-container-usage-instructions)
- [3.2 Local Usage Instructions](#32-local-usage-instructions)
- [4. API Overview](#4-api-overview)
- [5. Development Decisions](#5-development-decisions)
- [6. Pipeline Evaluation and Improvements](#6-pipeline-evaluation-and-improvements)
- [7. Future Improvements](#7-future-improvements)

---
## 1) Core Features 
⚙️ Modular components: ingest · chunker · embedder · vectordb · generator · evaluator

🧩 Config‑driven: YAML config for models, chunk sizes, overlap, DB endpoints

🚀 Dockerized: docker compose up to run Qdrant and the REST service

🔍 Deterministic evaluation: JSON test sets + metrics (Precision, Recall, F1, Cosine Similarity)

## 2) Tech Stack 

• Language: Go

• Vector DB: Qdrant (gRPC)

• Embeddings: Ollama/Nomic (local)

• Containerization: Docker & Docker Compose

• Testing: Go test

## 3.1) Docker Container Usage Instructions

To utilize the Rag Pipeline in a Docker container, follow these simple steps:
1. **Clone the Project**: 
  ```bash
    git clone https://github.com/ahoirg/rag-pipeline.git
    cd rag-pipeline
  ```

2. **Initialize Docker**: Ensure that Docker is running on your system. If not, start Docker from your system's applications.

3. **Build the Docker Image**:
  ```bash
    docker-compose up --build
  ```

4. **Send a Ping**: http://localhost:8080/api/ping
  >Expected response:
  ```bash
    {
      "success": true,
      "message": "Api is working...",
      "timestamp": "2025-11-12T23:27:43.2488977+01:00"
    }
  ```
  • Restful Api: http://localhost:8080/api/ 

## 3.2) Local Usage Instructions

**Note**: If you are trying the project in a Docker container, you should skip this instruction. Proceed directly to [4. API Overview](#4-api-overview) section.  

### Dependencies

| Name | Version |
|------|----------|
| **Go** | 1.24.10 |
| **Qdrant** | latest |
| **Ollama** | latest |
| **Embedding Model** (`nomic-embed-text`) | latest |
| (*)**Generator Models** (`llama3.2:3b`) | llama3.2:3b |

> (*) Optional models: `phi3:mini`, `tinyllama` or any other model available in your local Ollama installation.

### Prerequisites
  -All dependencies must be ready on the device
  - Qdrant running locally (gRPC on localhost:6334) 
  - Ollama running locally with models:
    - `ollama pull nomic-embed-text`
    -  `ollama pull llama3.2:3b`
    -  (optional) `ollama pull tinyllama` or `ollama pull phi3:mini`
    >Any non-default models you want to use at runtime must be configured in the config.yaml file.

1) **Run the application** : Use the following commands to download dependencies and start the server.
  ```bash
    go mod download
    go run main.go
  ```
⚠️ Local Setup Note: You need to set both the Qdrant and Ollama base URLs to localhost for the pipeline to run locally.
```yaml
qdrant:
  host: "localhost" <--
  port: 6334

ollama:
  base_url: "http://localhost:11434" <--
```

## 4) API Overview

| Method | Endpoint | Description |
|---------|-----------|-------------|
| **GET** | `/api/ping` | Health check endpoint |
| **GET** | `/api/evaluation/retrieval` | Returns retrieval evaluation results |
| **GET** | `/api/evaluation/generation` | Returns generation evaluation results |
| **POST** | `/api/storebook` | Stores a document into the vector database |
| **POST** | `/api/ask` | Full RAG workflow: retrieves relevant context and generates a final answer |
| **POST** | `/api/ask-directly` | Generates an answer directly without performing retrieval|

## 5) Development Decisions
We aimed to create a modular and flexible back-end and RAG pipeline. It  helps make it easier to implement future changes as the project grows.

### Backend
We used chi-go as the router because it makes managing API endpoints straightforward. It is fast, lightweight and provides useful middleware support. Chi also includes middleware for logging API requests and preventing crashes that may occur during request handling. Also it validates HTTP methods and returns appropriate status codes.

Folder Structure 📁
```yaml
├─ api/        # HTTP handlers (REST endpoints)
├─ db/         # Database layer
├─ eval_data/  # Labeled QA/eval JSON files ( for the evaluation of the rag papline, loaded automatically)
├─ evaluation/ # Evaluation logic (retrieval/generation metrics)
├─ models/     # Core data models (Chunk, Document, Embedding, etc.)
├─ services/   # Business logic: chunker, embedder, retriever, generator
├─ utils/      # Shared utilities 
├─ config.yaml # Configuration file
├─ main.go     # Application entrypoint
```

### RAG Pipeline

We use a main rag_service struct that manages the chunker, embedder, generator and communication with the vector database. This structure keeps each component modular. In the future, any part of the pipeline can be modified by simply updating the function bodies inside the corresponding service file, without affecting the rest of the system.

• Chunker: We implement word-based chunking using a sliding-window technique without relying on external frameworks. For sentence-aware chunking, we utilize existing Go libraries.

• Embedder: We support all embedding models that are based on Ollama. By default, we recommend using "nomic-embed-text", as it has a relatively small size and is ideal for the chunk lengths used in this project.

• Generator: We support all Ollama-based generator models. By default, we recommend "llama3.2:3b" (2GB, 128K context length), which easily handles our chunk token requirements. For a more lightweight option, TinyLlama (637MB) can be used, but it will fail when the number of chunks exceeds 4 due to its smaller context window.
> Ollama was chosen because it can be installed locally, requires no internet connection after initial setup and provides quick access to multiple models once integrated.

• Vector Database: Qdrant was chosen as the vector database because it can be easily integrated with Go and run locally. We use dense vector retrieval with cosine similarity. However, Qdrant also supports dense, sparse and hybrid search (multipvector) approaches.This flexibility allows us to quickly integrate other retrieval approaches into our system. [For more detail.](https://qdrant.tech/documentation/concepts/vectors/)

## 6) Pipeline Evaluation and Improvements
We prepared the evaluation data from the paragraph about the University of Notre Dame in the SQuAD 2.0 training set.

Retrieval and generation tests consist of 22 question-answer pairs. 
For retrieval tests, when the chunker changes, gold chunks are manually created with AI assistance for each chunking approach.

⚠️⚠️⚠️
>Each question–answer pair also contains information about the relevant chunks. Since these relevant chunks typically consist of only a single chunk, increasing the number of retrieved chunks (top-k) causes the F1 score to drop. Therefore, the recall metric is more important for our use case. Of course, to avoid incorrect or irrelevant information, precision and F1 should also be high.
>However, this limitation of the dataset can be handled through a customized prompt.

### First Evaluation (v0.0.1)
**Strategy**: split the text into chunks with sliding window approach --- **Chunk Size**: 300 words --- **Overlap** : 30 words (%10 overlap) --- **Top-K** : 2

<b>Retrieval Evaluation</b><br>
<table>
<tr>
<td>


| Metric | Value |
|--------|--------|
| Avg. Precision | **0.295** |
| Avg. Recall | **0.591** |
| Avg. F1 | **0.394** |

</td>
<td>

<b>Generation Evaluation</b><br>

| Metric | Value | Generator |
|--------|--------|-----------|
| Avg. Similarity Score (Cosine) | **0.513** | tinyllama |
| Avg. Similarity Score (Cosine) | **0.548** | llama3.2:3b |
| Avg. Similarity Score (Cosine) | **0.556** | phi3:mini |

</td>
</tr>
</table>

### Last Evaluation (v0.0.2)

**Strategy**: split the text into chunks with sliding window approach and ***customized prompt*** --- **Chunk Size**: 300 words --- **Overlap** : 55 words (%18 overlap) --- **Top-K** : 4

<table>
<tr>
<td>

<b>Retrieval Evaluation</b><br>

| Metric | Value |
|--------|--------|
| Avg. Precision | **0.215** |
| Avg. Recall | **0.636** |
| Avg. F1 | **0.315** |

</td>
<td>

<b>Generation Evaluation</b><br>

| Metric | Value | Generator |
|--------|--------|-----------|
| Avg. Similarity Score (Cosine) | **0.557** | tinyllama |
| Avg. Similarity Score (Cosine) | **0.797** | llama3.2:3b |
| Avg. Similarity Score (Cosine) | **0.842** | phi3:mini |

</td>
</tr>
</table>

### Our analyses:
We experimented with different chunking strategies such as Fixed-Length Chunking, Sentence-Based Chunking and Sliding Window Chunking. 
In approaches where we did not **use overlap** (i.e., chunking without any sliding window), we consistently observed poor retrieval performance.

For example, with Sentence-Based Chunking (chunk size: 5, overlap: 0, K = 2, using the new prompt), the results were:
``` json
{
    "AvgPrecision": 0.136363636363636,
    "AvgRecall": 0.25,   <-- important
    "AvgF1": 0.174242424242424
}
```
>The chunk and QA data used in the evaluation can be found under the path eval_data/previous_tests_data.zip/sources2

Therefore, in our experiments, the overlap value was always set to a value greater than 0.
In the Sentence-Based approach, the chunk sizes are not fixed. This situation can cause problems for models like TinyLLaMA and lead to exceeding the token limits expected by the model.
While generator scores in the sentence-based approach appeared high, manual inspection showed that these similarities were misleading. The generated answers often did not contain the actual answer. They included only similar vocabulary. 

Therefore we converted the texts into word slices and created chunks using a fixed number of words each time. We used cosine similarity for the retrieval step. In the future, we can make the retrieval better by using hybrid search or re-ranking methods. By changing the word overlap, chunk size, and top-k values, we achieved more stable and consistent results.

###  Comparison of  "v0.0.1" vs "v0.0.1+new promt" vs "v0.0.2"
| Version | Metric | Value | Generator | 
|--------|--------|-----------|---------|
|  v0.0.2 | Avg. Similarity Score (Cosine) | **0.556** | tinyllama     | 
|  v0.0.2 |Avg. Similarity Score (Cosine) | **0.798** ⬆️| llama3.2:3b  |
|  v0.0.2 |Avg. Similarity Score (Cosine) | **0.842** ⬆️| phi3:mini     |
| ────────────── | ────────────── | ───────── | ───────── |
|  v0.0.1+new promt | Avg. Similarity Score (Cosine) | **0.575** | tinyllama     | 
|  v0.0.1+new promt |Avg. Similarity Score (Cosine) | **0.533** | llama3.2:3b  |
|  v0.0.1+new promt |Avg. Similarity Score (Cosine) | **0.563** | phi3:mini     |
| ────────────── | ────────────── | ───────── | ───────── |
|  v0.0.1 | Avg. Similarity Score (Cosine) | **0.513** | tinyllama     | 
|  v0.0.1 |Avg. Similarity Score (Cosine) | **0.548** | llama3.2:3b  |
|  v0.0.1 |Avg. Similarity Score (Cosine) | **0.556** | phi3:mini     |


## 7) Future Improvements 
### Backend
1) **Testing**: Comprehensive unit tests, integration tests, end-to-ends should be added to improve reliability. Code coverage must be at least 80%.
2) **Configurable Parameters**: Runtime parameters (e.g.: top_k, temperature, embedding models, generator models) should be made fully configurable via API requests.
3) **CI/CD**: develop → main workflow with GitHub Actions. Merging to main triggers an automated build and a versioned release. It should create a package with version information.
4) **Managing Logging**: Structured JSON logging should be implemented to increase observability.
5) **File Architecture**: The project structure can be further organized to maintain clarity as the API continues to grow.  [For more detail.](https://medium.com/@smart_byte_labs/organize-like-a-pro-a-simple-guide-to-go-project-folder-structures-e85e9c1769c2)
   
### RAG Pipeline
1) **Different Retrieval Approaches:** Future improvements should include support for sparse vector retrieval and hybrid search methods.
2) **Benchmark Dataset:** A Gold Chunk test set should be created to evaluate different chunking strategies.
3) **Reranker integration:** Implement two-stage retrieval. For example, retrieve 10 chunks using cosine similarity, then rerank with cross-encoder to select top 3.
4) **Generation Optimization:** Different prompts should be tested to improve answer quality.




