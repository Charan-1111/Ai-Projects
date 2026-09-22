# AI Projects

A collection of small AI engineering projects built while learning how modern
LLM applications work. The projects are independent Go services, each with its
own `go.mod`, configuration, runtime dependencies, and API.

## Projects

| Project | Purpose | Main concepts |
| --- | --- | --- |
| [LLM Playground](LLM-Playground/README.md) | Generate text with Gemini, stream responses, discover configured models, and maintain chat conversations. | LLM APIs, streaming, retries, timeouts, token usage, pricing, persistence |
| [Semantic Search](Semantic-Search/README.md) | Create Gemini embeddings, compare text, store documents and chunks, and search with PostgreSQL/pgvector. | Embeddings, cosine similarity, chunking, vector databases, semantic search |

The repository is a learning workspace rather than a single deployable
application. Run each project from its own directory.

## Repository Layout

```text
Ai-Projects/
├── LLM-Playground/       # Gemini generation and conversation API
├── Semantic-Search/      # Gemini embeddings and pgvector search API
├── documentation/        # Learning notes, project status, and architecture notes
├── .gitignore
└── README.md
```

## Common Requirements

- Go 1.25 or later
- PostgreSQL
- A Google Gemini API key
- PowerShell, Bash, or another shell suitable for running Go commands

Each service loads a `.env` file from its current working directory. Local
environment files are ignored by Git. Create your own values; do not copy API
keys or database passwords into tracked files.

## Quick Start

### LLM Playground

The LLM Playground uses PostgreSQL for persistent chat conversations and
listens on `http://localhost:8000`.

```powershell
cd LLM-Playground
go mod download
```

Create `LLM-Playground/.env`:

```dotenv
LLM_PROVIDER_API_KEY=your-gemini-api-key
DB_USERNAME=postgres
DB_PASSWORD=your-database-password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=llm_playground
```

Start the service:

```powershell
go run ./cmd/llm
```

Useful endpoints:

- `GET /health`
- `GET /v1/llm/models/available`
- `POST /v1/llm/generate`
- `POST /v1/llm/generate/stream`
- `POST /v1/llm/chat`

The streaming endpoint uses Server-Sent Events. Chat requests can use the
`X-Session-Id` header to identify a conversation.

See [LLM-Playground/README.md](LLM-Playground/README.md) for request examples,
response formats, configuration details, and the persistence model.

### Semantic Search

Semantic Search uses PostgreSQL with the [pgvector](https://github.com/pgvector/pgvector)
extension and listens on `http://localhost:8000` by default.

```powershell
cd Semantic-Search
go mod download
```

Create `Semantic-Search/.env`:

```dotenv
GEMINI_API_KEY=your-gemini-api-key
DB_USERNAME=postgres
DB_PASSWORD=your-database-password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=semantic_search
```

The database user must be able to create the `vector` extension and the
application tables, or an administrator must install the extension first.

Start the service:

```powershell
go run ./cmd/llm
```

Useful endpoints:

- `GET /health`
- `POST /semantic/embed`
- `POST /semantic/compare`
- `POST /semantic/search`
- `POST /semantic/document/inject`
- `POST /semantic/document/inject/multiple`
- `PUT /semantic/document/:docId`
- `DELETE /semantic/document/:docId`
- `POST /semantic/chunks/upload/document`
- `POST /semantic/chunks/search`

The current indexing flow is synchronous: document embedding, word chunking,
chunk embedding, and database writes happen during the HTTP request. See
[Semantic-Search/README.md](Semantic-Search/README.md) for API payloads, schema
details, and current limitations.

## Configuration

Both services read `config/local/config.json` by default. Set
`CONFIG_FILE_PATH` to load a different configuration file:

```powershell
$env:CONFIG_FILE_PATH = "config/local/config.json"
go run ./cmd/llm
```

Configuration currently includes model definitions, retry values, SQL queries,
and service-specific settings. The files under `config/local/` are development
configuration, not production deployment manifests.

Important service differences:

- `LLM-Playground` currently listens on a hard-coded `:8000` address.
- `Semantic-Search` uses the JSON `port` setting and defaults to `:8000`.
- The Semantic Search provider currently uses the `gemini-embedding-001` model
	and 1536-dimensional vectors directly in code.

## Architecture At A Glance

```text
Client
	|
	+--> LLM Playground --> Gemini generation API
	|          |
	|          +--> PostgreSQL conversation history
	|
	+--> Semantic Search --> Gemini embedding API
						 |
						 +--> PostgreSQL + pgvector
										+--> documents
										+--> chunks
```

Both services use Fiber for HTTP routing, `pgx` for PostgreSQL access, and
request IDs for tracing responses. CORS is currently permissive and neither
service provides authentication or authorization.

## Development Checks

Run formatting, compilation, and tests from each project directory:

```powershell
gofmt -w .
go test ./...
go vet ./...
```

There are currently no automated tests in the repository, so a successful
build does not replace API-level verification against a configured Gemini key
and database.

## Documentation

- [Project structure notes](documentation/project-structure.md)
- [First-cut feature status](documentation/first-cut-feature-status.md)
- [Semantic Search learning roadmap](documentation/semantic-search.md)
- [Semantic Search indexing architecture](Semantic-Search/docs/indexing-architecture.md)

Some documentation describes planned learning milestones rather than shipped
behavior. The per-project READMEs are the most direct references for running
the current implementations.

## Current Limitations

- No authentication, authorization, or production deployment configuration.
- CORS allows all origins.
- No automated unit or integration test suite.
- Database schema and SQL are configuration-driven rather than managed by a
	migration tool.
- Semantic Search indexing is synchronous and sends one embedding request per
	chunk; background workers and batch embedding are planned, not implemented.
- Model names and pricing in local JSON configuration should be verified before
	using the services against a live provider account.

## Learning Direction

The projects progress from direct LLM integration to retrieval foundations:

```text
Gemini generation
		-> streaming, retries, usage, cost, conversations
		-> embeddings and similarity
		-> chunking and vector search
		-> background indexing and RAG evaluation
```
