# AI Projects

This repository contains a set of Go-based experiments focused on practical AI application patterns: LLM integration, document-grounded Q&A, semantic search, and tool calling. Each project is intentionally self-contained so it can be developed, run, and tested independently.

## Overview

The goal of this workspace is to explore building AI features with a hands-on, iterative approach. The projects here cover:

- LLM APIs and chat workflows
- Prompt-driven generation and streaming responses
- Semantic retrieval and embedding pipelines
- Vector search with PostgreSQL/pgvector
- Agent-like tool execution patterns
- Retrieval-augmented document Q&A

## Included projects

| Project | Focus | Main technologies |
| --- | --- | --- |
| [LLM-Playground](LLM-Playground/) | LLM generation, chat APIs, and service orchestration | Go HTTP services, Gemini providers, streaming responses, config-driven setup |
| [Semantic-Search](Semantic-Search/) | Embeddings, chunking, semantic similarity, and vector storage | Go, PostgreSQL, pgvector, cosine similarity, document indexing |
| [Ask-My-Documents](Ask-My-Documents/) | Document-centric Q&A and retrieval over local content | Retrieval workflow, context assembly, LLM prompting, document processing |
| [Tool-Calling](Tool-Calling/) | Function/tool-calling experiments for structured AI workflows | Tool registration, model-tool orchestration, request/response handling |

## Repository structure

```text
Ai-Projects/
├── Ask-My-Documents/      # Document-grounded RAG experiments
├── LLM-Playground/         # Chat and generation service examples
├── Semantic-Search/        # Embedding and similarity search service
├── Tool-Calling/           # Tool/function-calling prototypes
├── documentation/          # Notes, system docs, project planning docs
├── .gitignore
├── LICENSE
├── README.md
└── .github/
```

## Requirements

Most projects rely on a similar local setup:

- Go 1.25 or later
- PostgreSQL for vector or stateful services
- A valid Google Gemini API key
- Local project configuration values (often in `.env` or `config/*.json`)
- Shell access for running Go commands and service startup scripts

Do not commit API keys, database credentials, or other secret materials.

## Getting started

1. Clone the repository.
2. Open the project you want to explore.
3. Download dependencies:

```bash
go mod download
```

4. Create the required local configuration values for that project.
5. Run the service from that project directory with its project-specific entrypoint.

Example pattern:

```bash
cd LLM-Playground
go mod download
go run ./cmd/llm
```

If a project has its own README, follow that file for environment variables and API usage; it is usually the most accurate source for local setup.

## Suggested project order

If you are exploring the repo from scratch, this is a sensible order:

1. [LLM-Playground](LLM-Playground/) — understand API and model integration basics.
2. [Semantic-Search](Semantic-Search/) — learn embeddings, vectors, and retrieval patterns.
3. [Ask-My-Documents](Ask-My-Documents/) — see document-grounded Q&A in practice.
4. [Tool-Calling](Tool-Calling/) — extend the system into agent/tool orchestration.

## Documentation

The [architecture diagram](documentation/architecture.md) shows how the four services and their data stores and model/tool integrations fit together. The [documentation/](documentation/) folder also contains notes, feature status, and planning material for the workspace.

## Development checks

From an individual project directory, a typical validation flow is:

```bash
gofmt -w .
go test ./...
go vet ./...
```

This repository is primarily an experimental workspace, so validation is usually project-specific and should be paired with direct runtime checks against the configured API keys and dependencies.

## License

This repository is licensed under the MIT License. See [LICENSE](LICENSE) for details.
