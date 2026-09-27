# Ask My Documents — Code Review

I reviewed `Ask-My-Documents` on the `main` branch of your [Ai-Projects repository](https://github.com/Charan-1111/Ai-Projects).

The project has a good starting structure: a separate Go service calls your existing Semantic Search API, so you can reuse the chunking, embeddings, and pgvector work you have already built.

**Current status:** The endpoint retrieves relevant chunks and returns the Semantic Search response. It does not yet generate an answer with an LLM or provide citations, so it is not a complete RAG flow yet.

## Findings, in priority order

### 1. Upstream errors are ignored

In [`AskDocuments`](https://github.com/Charan-1111/Ai-Projects/blob/main/Ask-My-Documents/internal/services/askDocuments.go), errors from the Semantic Search API call and JSON unmarshalling are ignored. The service then returns `nil` as its error.

The [handler](https://github.com/Charan-1111/Ai-Projects/blob/main/Ask-My-Documents/internal/handlers/askDocument.go) also ignores the service error and returns `200 OK`.

**Why it matters:** If Semantic Search is unavailable or returns invalid JSON, your API can return an empty response that appears successful.

**Fix:** Return errors from the service and handle them in the HTTP handler.

### 2. The server can block after a listener error

In [`StartFiberServer`](https://github.com/Charan-1111/Ai-Projects/blob/main/Ask-My-Documents/internal/server/application.go), the code receives from `listenErr` in the `select`, then attempts to receive from the same channel again while logging. Only one error was sent, so the second receive blocks.

Capture the error on the first receive:

```go
case err := <-listenErr:
    app.log.Log.Error().Err(err).Msg("Server failed to listen")
```

### 3. Search is restricted to Go beginner documents

The [search request](https://github.com/Charan-1111/Ai-Projects/blob/main/Ask-My-Documents/internal/services/askDocuments.go) always sends:

```json
{
  "filters": {
    "category": "golang",
    "difficulty": "beginner"
  }
}
```

A relevant system design note, rate limiter note with another category, or advanced Go note will never be retrieved.

**Fix:** Make filters optional request fields. Search across all documents when no filters are supplied.

### 4. Configuration loading has errors

In [`config.go`](https://github.com/Charan-1111/Ai-Projects/blob/main/Ask-My-Documents/internal/config/config.go):

- The `Env` field uses the JSON tag `evn`, but the configuration file uses `env`.
- If `os.ReadFile` fails, the function still attempts to unmarshal the returned bytes.
- `sonic.Unmarshal(fileBytes, c)` is sufficient; `&c` adds an unnecessary pointer level.

### 5. The response type only represents search results

[`SemanticResponse`](https://github.com/Charan-1111/Ai-Projects/blob/main/Ask-My-Documents/internal/models/response.go) represents the Semantic Search API response. That is useful inside the service, but the final `/rag/ask/documents` response should contain an **answer** and **citations**.

For example:

```json
{
  "answer": "Your notes say the rate limiter uses a token bucket. [1]",
  "citations": [
    {
      "reference": 1,
      "document_id": "doc-123",
      "chunk_id": "chunk-456",
      "title": "Rate Limiter Notes"
    }
  ]
}
```

Also include `chunk_index` in your internal search result model because the Semantic Search API returns it.

### 6. Validate the question and limit LLM context

The current request model accepts an empty prompt. Validate it before calling Semantic Search.

When you add the LLM call, limit the **total size of the retrieved text** you send to the model. Retrieving five chunks does not guarantee that they fit comfortably in the model’s context window.

## Suggested implementation flow

Keep **Semantic Search** responsible for uploading documents, chunking, embeddings, storage, and retrieval.

Make **Ask My Documents** responsible for deciding whether retrieved chunks are useful, constructing the LLM context, generating an answer, and returning citations.

```text
User question
    ↓
Call Semantic Search
    ↓
Receive relevant chunks
    ↓
Are there useful chunks?
    ├── No → Say the documents do not contain enough information
    └── Yes
         ↓
       Build numbered context from chunks
         ↓
       Send question + context to LLM
         ↓
       Validate cited reference numbers
         ↓
       Return answer + citations
```

### Build it in this order

1. Fix error handling, configuration loading, and the listener bug.
2. Remove hardcoded search filters.
3. Test retrieval using questions whose answers are present and absent from your notes.
4. Build numbered excerpts containing the document title, document ID, chunk ID, and content.
5. Add an LLM call that answers using those excerpts.
6. Return structured citations and reject citation numbers that were never supplied.
7. Test failures: Semantic Search unavailable, invalid upstream JSON, no relevant chunks, and an invented citation.

**Important:** Checking that `[1]` refers to a supplied excerpt catches fabricated reference numbers. It does not prove that excerpt `[1]` actually supports the sentence. Review answer grounding separately when evaluating your RAG system.

## Overall assessment

The separate service is a good foundation, and reusing your Semantic Search API makes sense. The immediate priority is to make retrieval and error handling reliable. Then add context construction, LLM generation, and citations to complete the RAG flow.

This review is based on the repository source on `main`. I could not run `go test` because Go was not installed in the review environment.