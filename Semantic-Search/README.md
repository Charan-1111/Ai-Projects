# Semantic Search

Semantic Search is a Go HTTP service for generating Gemini embeddings, storing
documents and document chunks in PostgreSQL with pgvector, comparing text
similarity, and retrieving semantically similar documents.

The current implementation is a synchronous prototype. Document embedding,
chunking, chunk embedding, and database writes happen during the HTTP request.
The design for background indexing and batch embeddings is documented in
[`docs/indexing-architecture.md`](docs/indexing-architecture.md), but those
features are not implemented yet.

## Features

- Generate a 1536-dimensional embedding for arbitrary text.
- Calculate cosine similarity between two pieces of text.
- Store and update documents with title, content, category, source, metadata,
  and an embedding.
- Split content into overlapping word chunks and store chunk embeddings.
- Search whole documents by vector similarity.
- Search chunks by vector similarity with optional category, difficulty, and
  minimum-score filters.
- Create the `vector` PostgreSQL extension and application tables at startup.
- Add a request ID to responses through the `X-Request-ID` header.

## Requirements

- Go 1.25 or later
- PostgreSQL
- The PostgreSQL [pgvector](https://github.com/pgvector/pgvector) extension
- A valid Gemini API key with access to `gemini-embedding-001`

The application creates the pgvector extension and its tables at startup. The
database user therefore needs permission to create the extension and tables,
or the extension must already be installed by an administrator.

## Configuration

Run the application from the `Semantic-Search` directory. The entrypoint loads
`.env` from the current working directory, then reads
`config/local/config.json` unless `CONFIG_FILE_PATH` is set.

Create a local `.env` file with your own values:

```dotenv
GEMINI_API_KEY=your-gemini-api-key
DB_USERNAME=postgres
DB_PASSWORD=your-database-password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=semantic_search
```

Do not commit API keys or database passwords. If credentials have ever been
committed to a shared repository, revoke or rotate them.

The JSON configuration currently controls:

| Setting | Default | Purpose |
| --- | --- | --- |
| `port` | `8000` | HTTP port |
| `chunk_details.chunk_size` | `300` | Words per chunk |
| `chunk_details.chunk_overlap` | `50` | Overlapping words between chunks |
| `embedding_models` | Gemini embedding names | Descriptive model configuration |
| `retries` | See config file | Reserved retry configuration |
| `queries` | SQL in config | Table creation and repository queries |

The active provider currently uses `gemini-embedding-001` and output
dimensionality `1536` directly in code. The configured model and retry values
are not currently used by the provider.

## Database schema

On startup the service creates these tables if they do not already exist:

### `documents`

Stores the original document and its embedding.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `UUID` | Primary key |
| `title` | `TEXT` | Required |
| `content` | `TEXT` | Required |
| `category` | `TEXT` | Optional |
| `source` | `TEXT` | Optional |
| `metadata` | `JSONB` | Optional |
| `embedding` | `VECTOR(1536)` | Required |
| `indexing_status` | `TEXT` | Defaults to `pending` |
| `created_at`, `updated_at` | `TIMESTAMPTZ` | Database timestamps |

### `chunks`

Stores chunk content and embeddings. Chunks reference their parent document
with `ON DELETE CASCADE`.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `UUID` | Primary key |
| `document_id` | `UUID` | Required parent document |
| `content` | `TEXT` | Chunk text |
| `chunk_index` | `INTEGER` | Zero-based order |
| `start_position`, `end_position` | `INTEGER` | Word offsets |
| `embedding` | `VECTOR(1536)` | Required |
| `created_at` | `TIMESTAMPTZ` | Database timestamp |

The schema and SQL statements live in `config/local/config.json`; there is no
separate migration system yet.

## Run locally

From the project directory:

```powershell
cd Semantic-Search
go mod download
go run ./cmd/llm
```

The service listens on `http://localhost:8000` by default. Set `port` in the
JSON configuration to use another port. To use another configuration file:

```powershell
$env:CONFIG_FILE_PATH = "config/local/config.json"
go run ./cmd/llm
```

Verify that the server is running:

```powershell
Invoke-RestMethod http://localhost:8000/health
```

Expected response:

```json
{
  "code": 0,
  "message": "Server is healthy"
}
```

## HTTP API

All request bodies are JSON. Most successful responses use `code: 0`; body
binding failures return `400`, and provider/database failures generally return
`500`. There is currently no authentication or authorization.

### Generate an embedding

`POST /semantic/embed`

```json
{
  "text": "PostgreSQL stores relational data."
}
```

The response contains the original text and a 1536-value `embedding` array in
`contents`.

### Compare two texts

`POST /semantic/compare`

```json
{
  "original": "A quick brown fox",
  "compare": "A fast brown fox"
}
```

The response contains a numeric `similarity` value calculated with cosine
similarity.

### Store one document

`POST /semantic/document/inject`

```json
{
  "docTitle": "Vector databases",
  "docContent": "A vector database stores numerical representations of content.",
  "docCategory": "database",
  "docSource": "internal-notes",
  "metaData": {
    "difficulty": "beginner"
  }
}
```

The service generates a UUID, embeds `docTitle + " : " + docContent`, and
stores the document. The response returns the generated `docId` and
`embeddingStatus: "completed"`.

This endpoint stores only the parent document. To create searchable chunks,
use `/semantic/chunks/upload/document`.

### Store multiple documents

`POST /semantic/document/inject/multiple`

Send an array of document objects using the same fields as the single-document
endpoint. Documents are embedded and stored sequentially. If a later document
fails, earlier documents remain saved.

### Index documents into chunks

`POST /semantic/chunks/upload/document`

Send an array of document objects:

```json
[
  {
    "docTitle": "Vector databases",
    "docContent": "A vector database stores numerical representations of content.",
    "docCategory": "database",
    "docSource": "internal-notes",
    "metaData": {
      "difficulty": "beginner"
    }
  }
]
```

For each document, the service creates the parent document, splits its content
into 300-word chunks with 50 words of overlap by default, calls Gemini once per
chunk, and inserts the chunks in a transaction. A successful response means
the synchronous upload completed.

### Search whole documents

`POST /semantic/search`

```json
{
  "query": "How are vector representations stored?",
  "noOfDocs": 5
}
```

`noOfDocs` defaults to `10` when it is zero or negative. Results are ordered by
pgvector cosine distance and returned in `content.documentResults`.

The request model also accepts `filters` and `minimumScore`, but the current
whole-document search implementation ignores them.

### Search document chunks

`POST /semantic/chunks/search`

```json
{
  "query": "How are vector representations stored?",
  "noOfDocs": 5,
  "minimumScore": 0.7,
  "filters": {
    "category": "database",
    "difficulty": "beginner"
  }
}
```

Chunk search applies `category`, `difficulty`, and `minimumScore` when filters
are supplied. Provide a positive `noOfDocs`; unlike whole-document search,
chunk search does not currently apply a default limit.

### Update a document

`PUT /semantic/document/{docId}`

Send a document object. The service regenerates the document and chunk
embeddings, then transactionally updates the document, removes its old chunks,
and inserts the replacement chunk set. It returns `404` when the document does
not exist.

### Delete a document

`DELETE /semantic/document/{docId}`

Deletes the document and its chunks in a transaction. It returns `404` when the
document does not exist.

## Project structure

```text
cmd/llm/main.go                 Application entrypoint
config/local/config.json        Runtime configuration and SQL statements
docs/indexing-architecture.md  Current flow and planned indexing design
internal/chunks                 Word-based chunking
internal/config                 JSON configuration loading
internal/handlers               HTTP request handlers
internal/models                 Request and response types
internal/providers              Gemini embedding provider
internal/server                 Application construction and routes
internal/services               Embedding, indexing, and search workflows
internal/store/database         PostgreSQL and pgvector repository
```

## Current limitations

- Indexing is synchronous and can make one request perform many Gemini calls.
- There is no background queue, retry worker, batch embedding, or status API.
- There are no automated tests in the repository yet.
- There are no vector indexes beyond the table primary keys, so search
  performance may degrade as the dataset grows.
- `indexing_status` is not used to hide incomplete documents from search, and
  chunk uploads currently write the literal value `true` rather than a
  documented completed state.
- The application enables permissive CORS for all origins and has no auth.
- Several provider and database errors are returned to clients as generic
  messages, while some endpoints expose the underlying error text.
- Search response `durationMs` fields are present but are not populated.
- Configuration SQL is applied only at startup; schema upgrades require manual
  changes.

## Development checks

There are currently no test files. Run the compiler and package checks with:

```powershell
go test ./...
```

For a manual smoke test, start PostgreSQL and the service, call `/health`, then
exercise `/semantic/document/inject`, `/semantic/chunks/upload/document`, and
the two search endpoints with the examples above.

## Roadmap

The planned indexing work is tracked in
[`docs/indexing-architecture.md`](docs/indexing-architecture.md). It includes:

- background workers and observable `pending`/`processing`/`completed`/`failed`
  states;
- atomic replacement of chunk sets during reindexing;
- batch embedding with provider limits and cancellation;
- durable retry handling and document version checks; and
- status and reindex endpoints.

These items should be treated as design notes until the corresponding code and
API contracts are implemented.
