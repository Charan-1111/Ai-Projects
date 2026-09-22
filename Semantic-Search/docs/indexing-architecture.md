# Semantic Search Indexing Architecture

This document defines the indexing behavior for documents and their chunks in
`Semantic-Search`. It records the current synchronous behavior and the planned
background indexing and batch embedding changes.

## Current implementation

The current chunk upload flow is synchronous:

```text
POST /semantic/chunks/upload/document
        |
        +-- create and save the document embedding
        +-- split the document into chunks
        +-- generate one embedding request per chunk
        +-- save all chunks in a database transaction
        +-- mark the document as indexed
        |
        +-- return after indexing finishes
```

The current provider interface exposes `EmbedText`, which accepts one string at
a time. Therefore a document with 100 chunks currently causes 100 embedding
requests. The database schema already stores an `indexing_status` on documents
and uses `pending` as its default value. The chunk table references the parent
document with `ON DELETE CASCADE`.

The document update path must replace the complete chunk set. Updating only the
document embedding would leave stale chunk content and stale chunk vectors in
search results. The replacement of the document and its chunks must remain
atomic.

## Design goals

The indexing design must:

- return quickly after accepting a document;
- make indexing status observable;
- never expose a partially indexed document through chunk search;
- replace stale chunks when document content changes;
- make retries safe and idempotent;
- preserve the order of input texts and output embeddings in batch requests;
- respect provider input limits, token limits, rate limits, and timeouts;
- retain the synchronous flow as a simple development and fallback mode.

## Indexing state machine

Use these states for a document:

```text
pending -> processing -> completed
                    \\-> failed
```

State meanings:

| State | Meaning |
| --- | --- |
| `pending` | The document is saved, but chunk indexing has not started. |
| `processing` | A worker has claimed the document and is indexing it. |
| `completed` | The document and its complete current chunk set are searchable. |
| `failed` | The latest indexing attempt failed. The failure reason and attempt count should be recorded. |

Only `completed` documents should be returned by chunk search. A failed update
must not delete the last known-good chunk set unless the product explicitly
chooses to make the document unavailable during reindexing.

## Target background flow

```text
Client
  |
  +-- POST /semantic/document/inject
  |       |
  |       +-- validate request
  |       +-- create document with indexing_status = pending
  |       +-- enqueue document ID
  |       +-- return document ID and pending status
  |
  +-- worker claims document
          |
          +-- atomically change pending -> processing
          +-- split current document content into chunks
          +-- generate chunk embeddings
          +-- transactionally replace the chunk set
          +-- mark document completed
          |
          +-- on failure: mark document failed and store error details
```

The HTTP handler must not perform chunking or call the embedding provider after
it has accepted the document. The worker owns all long-running indexing work.

### Worker requirements

The first implementation can use an in-process worker to keep the learning
surface small:

- start a fixed number of workers during application startup;
- use a bounded channel as the in-memory queue;
- enqueue only the document ID, not the full document payload;
- read the latest document row after claiming the job;
- stop workers gracefully during application shutdown;
- use a context with a timeout for each provider request and indexing job.

An in-process queue is not durable. Jobs can be lost if the process stops before
they complete. A database-backed job table or external queue is required when
job durability, multiple application instances, or guaranteed retry delivery
is needed.

### Claiming and concurrency

A worker must claim a document conditionally so two workers cannot index the
same version at the same time. The claim should be equivalent to:

```sql
UPDATE documents
SET indexing_status = 'processing', updated_at = NOW()
WHERE id = $1 AND indexing_status IN ('pending', 'failed')
RETURNING id;
```

For updates, add a document version or generation number. The worker captures
the version when it starts and commits results only when that version is still
current. This prevents an older, slower indexing job from overwriting chunks
for a newer update.

A delete must cancel or invalidate queued work. The final database transaction
must verify that the document still exists and has the expected version before
publishing chunks.

## API contract

### Create document

`POST /semantic/document/inject`

The request creates the document and queues indexing. It should return `202
Accepted`:

```json
{
  "document_id": "doc_123",
  "indexing_status": "pending"
}
```

The response must not claim `completed` before the worker has committed the
chunks.

### Get document status

Add a status endpoint, for example:

`GET /semantic/document/:docId/status`

Example response:

```json
{
  "document_id": "doc_123",
  "indexing_status": "processing",
  "indexing_attempt": 1,
  "error": null,
  "updated_at": "2026-09-22T12:00:00Z"
}
```

Use `200 OK` for an existing document, `404 Not Found` for an unknown document,
and expose a sanitized failure message rather than provider credentials or raw
internal errors.

### Retry indexing

A retry endpoint may requeue a failed document:

`POST /semantic/document/:docId/reindex`

It should return `202 Accepted` and be idempotent while the document is already
`pending` or `processing`.

## Transaction boundary

Embedding generation happens before the publish transaction. The worker should
then use one database transaction to:

1. verify the document version is still current;
2. delete the old chunks for that document;
3. insert the complete newly generated chunk set;
4. set `indexing_status = 'completed'`;
5. record the indexed version and completion time;
6. commit.

If any insert fails, roll back the transaction. The old chunk set remains
available, or the document remains excluded from search according to the
chosen visibility policy. Never publish only a subset of a document's chunks.

Document deletion should also remain transactional. Delete its chunks and the
parent document together, and invalidate any queued indexing job by checking the
parent row/version before publishing.

## Batch embedding

### Provider interface

Extend the provider abstraction with a batch operation while retaining the
single-text method for compatibility:

```go
type LLMProvider interface {
    EmbedText(ctx context.Context, text string) ([]float32, error)
    EmbedTexts(ctx context.Context, texts []string) ([][]float32, error)
}
```

`EmbedTexts` must satisfy these rules:

- reject an empty input list;
- return one vector per input and preserve input order;
- return an error if the provider returns a different number of vectors;
- use the configured embedding model and output dimensionality;
- honor the supplied context for cancellation and timeout;
- return no partial result for a failed batch unless partial results are an
  explicit, documented provider contract.

The Gemini implementation can pass multiple `Content` values to
`Models.EmbedContent`. The provider adapter, rather than the service, should
translate SDK responses into `[][]float32` and validate the response length.

### Batching policy

The indexing service should split chunks into bounded batches:

```text
chunks 1-20   -> batch request 1
chunks 21-40  -> batch request 2
chunks 41-60  -> batch request 3
```

The batch size must be configuration, not a hard-coded provider assumption. A
batch is valid only when it satisfies all configured limits:

- maximum input count;
- maximum estimated tokens or characters;
- provider request payload size;
- request timeout;
- rate-limit and concurrency limits.

If a document's chunks exceed the token or payload limit, split the batch
further. A simple first implementation can use a maximum input count and a
conservative character estimate, then replace it with the provider's tokenizer
when exact token accounting is required.

Suggested configuration:

```json
{
  "embedding": {
    "batch_size": 20,
    "max_input_tokens": 0,
    "max_concurrent_batches": 2,
    "request_timeout_ms": 30000
  }
}
```

`max_input_tokens: 0` can mean that token-based splitting is disabled and only
the provider's input-count and payload safeguards are applied. The actual
limits must be verified against the selected provider and model before enabling
large batches.

### Batch failure and retry behavior

Treat each batch as an independently retryable operation, but publish nothing
until every batch for the document succeeds. Retry transient provider failures
with bounded exponential backoff and jitter. Do not retry invalid requests,
unsupported dimensions, or authentication failures without changing the
configuration or credentials.

A failed batch causes the indexing job to fail. The worker may keep vectors in
memory for completed batches, but it must discard them if the document job is
aborted. The next retry should regenerate or safely reuse only results whose
request identity and document version are known.

## Recommended implementation order

1. Keep the current synchronous flow as the baseline and make its status values
   accurate.
2. Add `EmbedTexts` to the provider interface and Gemini adapter.
3. Add a batching helper with count, payload, timeout, and concurrency limits.
4. Change synchronous chunk indexing to use the batch helper and verify that all
   vectors are present before publishing.
5. Add document versioning and a status response model.
6. Add an in-process bounded worker queue and return `202 Accepted` from create.
7. Add status and reindex endpoints.
8. Add durable job storage if deployments require restart recovery or multiple
   API instances.
9. Add integration tests for update, delete, retry, stale worker results, and
   partial batch failure.

## Required tests

- A document with zero or invalid content is rejected without creating a job.
- A successful batch returns vectors in the same order as its input texts.
- A provider response with the wrong vector count fails the batch.
- A transient batch failure retries within the configured limit.
- A permanent batch failure marks the document `failed`.
- A successful update replaces all old chunks and makes only the new content
  searchable.
- An older indexing job cannot overwrite a newer document version.
- Deleting a document prevents a queued job from publishing chunks.
- A failed indexing transaction leaves no partial chunk set.
- Chunk search excludes documents that are not `completed`.

## Query-embedding cache

Query embedding is on the read path, so it can be cached after basic search
behavior is correct and covered by tests. Equivalent normalized queries should
reuse the same vector for the same embedding configuration:

```text
"How do goroutines work?"
"how do goroutines work?"
        |
        +-- normalize -> "how do goroutines work?"
        +-- cache key -> embedding-model-v1:how do goroutines work?
```

The cache key must include every setting that can change the vector space. At a
minimum include:

- embedding provider and model name;
- model or embedding configuration version;
- task type, if the provider supports task-specific embeddings;
- output dimensionality;
- normalized query text.

For example:

```text
gemini:gemini-embedding-001:v1:retrieval:1536:how do goroutines work?
```

Never use only the query text as the key. Vectors from different models,
versions, task types, or dimensions must not be mixed.

### Normalization

Normalization must be deterministic and conservative. The initial normalizer
should:

1. trim leading and trailing whitespace;
2. collapse consecutive whitespace to one space;
3. apply Unicode case folding or lowercase consistently with the product's
   language requirements.

Do not remove punctuation, reorder words, or apply stemming unless search
requirements explicitly define those transformations. Over-normalization can
make distinct queries share a vector unexpectedly. Do not cache an empty query.

### Cache behavior

Use a bounded cache with a TTL. An in-process LRU cache is sufficient for a
single application instance and is easy to add around the existing query
embedding call. A shared cache such as Redis is needed when multiple API
instances must share entries.

Required behavior:

- check the cache after normalization and before calling the provider;
- store only successful embedding responses;
- never cache provider errors or cancellation results;
- enforce a maximum entry count and/or maximum memory budget;
- expire entries after a configured TTL;
- invalidate or namespace entries when the model/version/configuration changes;
- copy vectors on return or otherwise prevent callers from mutating cached
  values;
- use request coalescing for identical concurrent misses when provider traffic
  needs to be limited.

Suggested configuration:

```json
{
  "query_embedding_cache": {
    "enabled": false,
    "ttl_seconds": 3600,
    "max_entries": 10000
  }
}
```

The cache should be disabled by default until search correctness is verified.
Cache hits and misses should be observable through counters and latency
metrics, without logging the complete query when queries may contain sensitive
content.

### Cache tests

- Equivalent normalized queries produce one provider call and two cache hits
  after the first request.
- Different model versions never share an entry.
- Different task types or output dimensions never share an entry.
- Failed provider calls are not cached.
- Expired entries call the provider again.
- The cache evicts entries at its configured bound.
- Concurrent identical misses do not create an avoidable provider-request
  stampede when request coalescing is enabled.

## Evaluating search quality

A technically correct API can still return poor search results. Maintain a
small, versioned evaluation set and run it whenever chunking, embeddings,
thresholds, or ranking logic changes.

### Initial evaluation set

| Query | Expected document or concept |
| --- | --- |
| How does Go run lightweight concurrent tasks? | Goroutines |
| How do goroutines communicate? | Channels |
| How can I cancel a request? | Context cancellation |
| Why do indexes make SQL queries faster? | Database indexing |

The expected value should identify a document or chunk that is known to contain
the answer. Keep the evaluation documents fixed while comparing configurations.
Add more queries over time, including paraphrases, ambiguous wording, and
queries that should return no result.

### Metrics

For each query record:

- `Hit@1`: the expected result is ranked first;
- `Hit@3`: the expected result appears in the first three results;
- similarity score of the expected result;
- end-to-end search latency in milliseconds;
- optionally, provider latency and cache hit/miss status.

Calculate aggregate hit rates as:

```text
Hit@1 = queries with expected result at rank 1 / total queries
Hit@3 = queries with expected result in ranks 1-3 / total queries
```

For example, if 8 of 10 queries return the expected result first:

```text
Hit@1 = 8 / 10 = 80%
```

Report the mean and a high percentile such as p95 for latency. Similarity
scores are useful for diagnosing thresholds, but they are not directly
comparable across different embedding models unless the scoring behavior is
known to be equivalent.

### Evaluation record

Store one record per query and configuration, for example:

```json
{
  "evaluation_version": "2026-09-22",
  "configuration": {
    "embedding_model": "gemini-embedding-001",
    "chunk_size": 300,
    "chunk_overlap": 50,
    "minimum_score": 0.7,
    "search_mode": "semantic"
  },
  "query": "How do goroutines communicate?",
  "expected_document": "Channels",
  "returned_document_ids": ["doc_12", "doc_07", "doc_31"],
  "expected_rank": 1,
  "expected_similarity": 0.91,
  "latency_ms": 42
}
```

The evaluation runner should use the same search endpoint and request shape as
normal clients. Record the dataset version, corpus version, model version, and
configuration with every run so results remain reproducible.

### Controlled comparisons

Change one variable at a time where possible, rebuild the index when chunk size
or overlap changes, and run the complete evaluation set for each configuration.
Compare:

- chunk size 300 versus 500;
- overlap 50 versus 100;
- different embedding models or model versions;
- different minimum similarity thresholds;
- pure semantic search versus hybrid lexical and semantic search.

Do not select a configuration using `Hit@1` alone. Prefer a configuration that
improves `Hit@1` and `Hit@3` without unacceptable latency or a large increase
in false-positive results. Keep baseline results so regressions are visible.

### Evaluation tests and safeguards

- Run the evaluation set against a fixed fixture corpus in CI or a repeatable
  local command.
- Fail or warn when a metric falls below the accepted baseline.
- Verify that expected documents are not excluded by status or minimum-score
  filtering.
- Include latency measurements separately from provider calls so query-cache
  improvements can be evaluated without hiding ranking regressions.

## Hybrid search

Vector search is strong at matching meaning, while keyword search is strong at
matching exact text. For example:

```text
Semantic query: lightweight concurrent tasks
Semantic match: goroutines

Keyword query: DPI-1047
Keyword match: a document containing DPI-1047
```

Hybrid search combines both signals. It is especially useful for error codes,
function names, product names, IDs, technical acronyms, and exact API names.

Do not introduce hybrid ranking until basic vector search is working and the
evaluation set has established a baseline. Otherwise it is difficult to tell
whether a result change comes from a better ranking strategy or from an
incorrect vector-search implementation.

### Retrieval flow

```text
Query
  |
  +-- vector retrieval: top K semantic candidates
  +-- keyword retrieval: top K lexical candidates
  +-- merge candidates by document/chunk ID
  +-- normalize the two scores
  +-- calculate the final score
  +-- rank and return the requested number of results
```

Run both retrieval paths with a larger candidate limit than the final response
limit. For example, retrieve 50 candidates from each path before merging the
top 10. The exact limits should be measured because larger candidate sets
increase database work and latency.

### Score fusion

The two retrieval systems usually produce scores on different scales, so do not
add raw scores directly. Normalize each score within its candidate set or use a
ranking-based fusion method. A weighted score can then be defined as:

```text
final_score = vector_weight * normalized_vector_score
            + keyword_weight * normalized_keyword_score
```

The weights must be configuration, for example:

```json
{
  "search": {
    "mode": "semantic",
    "vector_weight": 0.7,
    "keyword_weight": 0.3,
    "candidate_limit": 50
  }
}
```

Start with semantic search as the default mode. Make `hybrid` an explicit mode
while it is being evaluated. Do not assume that a keyword match should always
win: an exact match in a low-quality or unrelated field may be less useful than
a strong semantic match.

### Keyword retrieval options

For PostgreSQL, begin with a lexical query that supports the fields users need
to search exactly. Options include:

- PostgreSQL full-text search with a `tsvector` column and `tsquery` input;
- PostgreSQL trigram matching for identifiers, partial names, and typo-tolerant
  matching;
- exact equality or prefix matching for structured IDs and error codes.

Choose fields deliberately. Content, title, source, and selected metadata may
need different weights. Preserve punctuation and case-sensitive forms where
identifiers require it, even if the semantic query normalizer lowercases the
embedding input.

### Exact-match safeguards

Exact identifiers should be protected from being diluted by score fusion. If a
query contains a recognized error code, function name, or ID and a candidate
contains the exact token in a searchable field, apply a documented exact-match
boost or include that candidate in a guaranteed candidate set. The boost must
be bounded so unrelated exact matches cannot dominate every result.

Do not log full queries or document content when technical identifiers may be
sensitive. Record only the metrics and identifiers needed for evaluation.

### Evaluation requirements

Add exact-match cases to the evaluation set, such as an error code, a function
name, and an API name. Compare semantic and hybrid modes using the same corpus
and queries, measuring:

- `Hit@1` and `Hit@3`;
- exact-match recall for identifier queries;
- false-positive results caused by keyword matches;
- mean and p95 search latency;
- vector and keyword candidate counts.

Tune score normalization and weights against the evaluation set, changing one
variable at a time. Keep the pure vector baseline and reject hybrid search if
it improves exact-match queries but causes unacceptable regressions on ordinary
natural-language queries.

## Reranking

Vector search is fast enough to retrieve potential matches, but its ordering is
not always perfect. A reranking stage examines a small candidate set more
carefully before returning the final results:

```text
Query
  |
  +-- vector search retrieves 20 candidates
  +-- reranking model scores each candidate against the query
  +-- return the best 5 candidates
```

Vector search answers:

> Which records might be relevant?

The reranker answers:

> Among these candidates, which ones best answer this exact query?

Reranking is a later-phase optimization. Implement and evaluate basic vector
search first, then hybrid search if exact-match behavior requires it. A
reranker adds provider latency, cost, failure modes, and operational
complexity, so it should not be placed on every request by default without
evidence that its quality improvement justifies those costs.

### Reranking contract

The search service should:

1. retrieve a bounded candidate set using vector or hybrid search;
2. send the query and candidate text, title, and selected metadata to the
   reranking provider;
3. validate that every candidate receives a score or ranking;
4. sort candidates by reranker score while preserving a stable tie-breaker;
5. return only the requested result limit.

The candidate limit must be greater than the final result limit. For example,
retrieve 20 candidates and return the best 5. Measure larger candidate limits
because reranking cost and latency generally grow with the number and size of
candidate inputs.

### Provider interface

Keep reranking behind an abstraction so the search service is independent of a
specific model provider:

```go
type Reranker interface {
    Rerank(ctx context.Context, query string, candidates []Candidate) ([]float32, error)
}
```

The result must preserve candidate identity or return scores in exactly the
same order as the input candidates. Reject a response with missing, duplicate,
or extra scores. Use a request timeout and bounded candidate text length to
avoid unexpectedly large provider requests.

### Failure and fallback behavior

Reranking must not make search unavailable when the initial retrieval succeeded.
For a timeout, rate limit, provider error, or invalid response, record the
failure and return the original vector or hybrid ranking. Do not silently
return an empty result set. The fallback policy should be observable through
metrics and optionally exposed in debug-only response metadata.

Do not retry reranking indefinitely. Use one bounded retry only for known
transient errors, and honor the request context deadline. If the request does
not have enough remaining time, skip reranking and return the initial ranking.

### Evaluation and rollout

Add reranking as an explicit search mode or feature flag. Compare it with the
same fixed corpus and evaluation set used for vector and hybrid search. Measure:

- `Hit@1` and `Hit@3` improvement;
- answer or document relevance for the final result limit;
- reranking latency and end-to-end p50/p95 latency;
- provider cost and error rate;
- fallback frequency;
- quality impact for both semantic questions and exact identifier queries.

Start with a small candidate limit and offline evaluation. Enable it gradually
only when the relevance improvement is consistent and the added latency and
cost fit the product requirements. Keep the non-reranked ranking as the
default fallback and regression baseline.
