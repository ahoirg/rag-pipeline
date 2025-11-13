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
- [6. Pipeline Evaluation & Improvements](#6-Pipeline-Evaluation-&-Improvements)
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

**Note**: If you are trying the project in a Docker container, you should skip this instruction. Proceed directly to [2. Using the Project](#2-using-the-project) section.  

### Dependencies

| Name | Version |
|------|----------|
| **Go** | 1.24.10 |
| **Qdrant** | latest |
| **Ollama** | latest |
| **Embedding Model** (`nomic-embed-text`) | latest |
| (*)**Generator Models** (`tinyllama`) | TinyLlama-1.1B |

> (*) Optional models: `llama3.2:3b`, `phi3:mini`, or any other model available in your local Ollama installation.

### Prerequisites
  -All dependencies must be ready on the device
  - Qdrant running locally (gRPC on localhost:6334) 
  - Ollama running locally with models:
    - `ollama pull nomic-embed-text`
    -  `ollama pull tinyllama`
    -  (optional) `ollama pull llama3.2:3b` or `ollama pull phi3:mini`
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
## 6) Pipeline Evaluation & Improvements
### First Evaluation
Strategy:
<table>
<tr>
<td>

<b>Retrieval Evaluation</b><br>

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

### Last Evaluation

Strategy:

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




