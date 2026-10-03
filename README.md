# AI Projects

This repository is a collection of Go-based AI experiments and learning projects built while exploring modern LLM workflows, vector search, document Q&A, and tool-calling patterns. Each project is designed to be run independently from its own directory.

## Current Projects

| Project | Purpose | Main concepts |
| --- | --- | --- |
| [LLM-Playground](LLM-Playground/) | Lightweight LLM service for text generation, chat, and streaming responses. | Gemini APIs, streaming, chat persistence, retry handling, configuration |
| [Semantic-Search](Semantic-Search/) | Embedding and similarity search experiments with PostgreSQL and pgvector. | Embeddings, cosine similarity, chunking, vector search, database indexing |
| [Ask-My-Documents](Ask-My-Documents/) | Document-grounded Q&A and retrieval experiments over local content. | Retrieval, document chunking, contextual prompting, question answering |
| [Tool-Calling](Tool-Calling/) | Function/tool-calling prototypes for structured agent workflows. | Tool schemas, function orchestration, model/tool integration |

The repository is a learning workspace rather than a single deployable application. Most projects are intentionally modular and can be developed, run, and debugged separately.

## Repository Layout

```text
Ai-Projects/
├── Ask-My-Documents/      # Document-grounded Q&A experiments
├── LLM-Playground/         # LLM generation and chat APIs
├── Semantic-Search/        # Embedding and semantic search service
├── Tool-Calling/           # Tool/function-calling examples
├── documentation/          # Notes, architecture docs, and learning material
├── .gitignore
├── LICENSE
├── README.md
└── .github/
```

## Requirements

Most projects in this repository use the following workflow:

- Go 1.25 or later
- PostgreSQL for services that store state or vector data
- A Google Gemini API key
- A local `.env` file in the project directory when required
- PowerShell, Bash, or another shell for running Go commands

Each project keeps its own configuration and environment values locally. Do not commit API keys or database credentials.

## Getting Started

1. Clone the repository.
2. Move into a project directory of interest.
3. Download dependencies:

```bash
go mod download
```

4. Create a local `.env` file (if the project expects one) using the sample values from that project.
5. Run the service or app from that folder with its project-specific entrypoint.

Some projects include their own README or configuration notes inside their directory. For the most complete usage examples, inspect the project folder you want to run.

## Project Notes

### LLM Playground

The LLM Playground is the most mature project in the repo and includes chat and generation APIs, environment configuration, and service startup examples.

### Semantic Search

Semantic Search focuses on vector embeddings, document chunking, PostgreSQL storage, and similarity search workflows using pgvector.

### Ask My Documents

This project is aimed at document-grounded AI experiences, making it suitable for retrieval and question-answering experiments over local or uploaded content.

### Tool Calling

This project explores function and tool calling patterns, where an LLM can decide when to invoke external tools or structured operations.

## Development Checks

Run formatting and basic validation from each project directory when working on an individual service:

```bash
gofmt -w .
go test ./...
go vet ./...
```

There are currently no repository-wide automated tests, so practical verification should include running the relevant service and validating its output against your configured API keys and databases.

## Documentation

The `documentation/` folder contains notes, architecture summaries, and planning material related to the projects in this repository.

## License

This repository is licensed under the MIT License. See [LICENSE](LICENSE) for details.
