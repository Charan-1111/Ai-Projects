# Repository Assistant — Phase 1

## Goal

Build a Go pipeline that reads a specific commit of your Rate Limiter repository, extracts searchable code chunks, stores embeddings, and answers one repository question with verified source citations.

**End-of-phase demo:** Ask “Where is the limiter created?” through LLM Playground and receive an answer linking to the relevant declaration and initialization call at the indexed commit.

Plan for six morning sessions. Each session has one observable deliverable.

## Implementation sequence

| Session | Main deliverable | Done when |
|---|---|---|
| 1 | Read a pinned Git snapshot and print Go declarations | Output includes commit, file, physical line range, and symbol. Uncommitted edits leave the output unchanged. |
| 2 | Create Go-aware chunks | Functions, methods, types, and large declaration splits map to the exact committed source. |
| 3 | Persist chunks and embeddings | Indexing the same inputs creates no duplicate ready generation. Failed ingestion stays hidden from search. |
| 4 | Combine exact and semantic retrieval | Both a constructor name and a natural-language question retrieve the expected evidence. |
| 5 | Generate answers with citations | LLM Playground answers using retrieved evidence. Invalid evidence IDs cannot produce links. |
| 6 | Evaluate retrieval and answers | You can tell whether a failed answer came from missing evidence or incorrect interpretation. |

## Session 1 — Read committed code

Create a Repository-Assistant directory inside AI Projects. Begin with a command that accepts a local checkout and revision. The following command describes the intended interface to implement.

```bash
go run ./cmd/repository-assistant inspect \
  --repo /path/to/Rate-Limiter \
  --revision HEAD
```

1. Resolve the supplied revision once to a full commit SHA.
2. Enumerate the files in that commit using Git tree objects.
3. Select regular Go source files within configured file and size limits. Skip symlinks, submodules, vendor directories, generated files, and denied paths with recorded reasons.
4. Read each selected blob from Git using its object ID. Working-tree files may differ from the selected commit.
5. Parse the committed bytes and extract functions, methods, and type declarations.
6. Print the commit, path, line range, kind, and symbol.

Use Go's process API with explicit argument arrays for Git commands. Keep the local repository path supplied by the operator. Validate the revision, resolve it to an object ID, and use that ID for subsequent operations. Parse NUL-delimited tree output so filenames containing spaces remain intact.

Suggested in-memory representation

```go
type Declaration struct {
    CommitSHA string
    Path      string
    Package   string
    Kind      string
    Name      string
    Receiver  string
    StartLine int
    EndLine   int
    Source    string
}
```

Use one-based, inclusive physical line ranges. Retain the original source bytes for each declaration. Go AST end positions are exclusive, so calculate the end line from the final included byte. Use physical positions so line directives cannot change citation locations.

Illustrative output

```text
Commit: <full-commit-sha>

internal/limiter/token_bucket.go:12-19  type TokenBucket
internal/limiter/token_bucket.go:21-32  func NewTokenBucket
internal/limiter/token_bucket.go:34-57  method (*TokenBucket).Allow
```

The example filenames and symbols are placeholders. Use the actual declarations from your repository.

### Session 1 completion checks

- [ ] One function, method, and type match the committed file and printed lines.
- [ ] A harmless uncommitted edit leaves output unchanged for the same SHA.
- [ ] A newly created untracked file is absent from the selected snapshot.
- [ ] Parse failures and skipped files are visible in the output.

## Session 2 — Create chunks

Start with one chunk per function, method, or type declaration. Include useful declaration comments when appropriate and preserve their actual source range. Split unusually large declarations at line boundaries and retain the parent symbol on every part. Use a configurable size budget that fits your embedding provider.

Keep two text fields conceptually separate.

| Field | Purpose |
|---|---|
| Source text | Exact committed bytes corresponding to the cited line range |
| Search text | Source plus package, symbol, receiver, or repeated signature to improve retrieval |

Record repository ID, commit SHA, path, package, symbol, receiver, chunk ordinal, start line, end line, and both text fields. Include selected SQL and non-sensitive config when needed to explain the code. Use bounded ranges for these files and report unsupported formats.

### Session 2 completion checks

- [ ] Every chunk can be located in its committed file.
- [ ] A large function produces multiple chunks with valid original ranges.
- [ ] Search metadata added to a chunk never changes its citation range.

## Session 3 — Store embeddings

Reuse the existing embedding client and PostgreSQL + pgvector setup. Define an index generation using repository ID, commit SHA, chunker version, and embedding profile. The profile identifies the model, vector dimension, and any embedding settings that affect compatibility.

1. Create or claim an indexing generation.
2. Store chunks and their metadata with unique keys over generation, path, and ordinal.
3. Generate and validate embeddings, including vector dimension.
4. Persist successful batches while the generation remains unavailable to search.
5. Mark the generation ready only after all required chunks have embeddings.

For the initial CLI, allow one writer per generation. Make retries idempotent through unique keys and upserts. A later concurrent service needs leases or equivalent ownership controls. Preserve existing ready generations when a new generation fails. Record skipped files independently from failed embedding batches.

### Session 3 completion checks

- [ ] Reindexing the same inputs creates no duplicate ready generation or chunks.
- [ ] A simulated embedding failure exposes no partial generation to search.
- [ ] An earlier ready generation remains queryable after a failed rebuild.
- [ ] Incompatible vector dimensions are rejected.

## Session 4 — Retrieve relevant code

Implement a search function scoped to one ready generation. Combine exact symbol and literal matching with semantic vector retrieval. Use the same embedding profile for the question and indexed chunks.

Prioritize qualified symbol matches when the question contains an exact symbol. For natural-language questions, merge ranked exact and semantic candidates, deduplicate overlapping ranges, and cap repeated results from one file. Log each retrieval branch's candidates and final ranks.

Try these queries using actual symbols from your repository.

- “Where is the limiter created?”
- The exact constructor name.
- “Where does the application check whether a request is allowed?”

### Session 4 completion checks

- [ ] Expected source appears among the top results for both natural-language and exact-symbol queries.
- [ ] Results include repository, commit, path, lines, source, and rank metadata.
- [ ] Search never returns chunks from another generation.

## Session 5 — Answer with verified citations

Send the retrieved excerpts to LLM Playground with evidence IDs. Ask the model to return code claims and supporting evidence IDs. Allow an incomplete answer when the snippets do not establish the requested behavior.

Register evidence IDs only for excerpts actually delivered to the model. In Go, validate those IDs and construct source links using a trusted GitHub origin, the full commit SHA, the escaped file path, and verified line numbers. For local fixtures without a GitHub origin, render path, commit, and line references.

The answer to “Where is the limiter created?” should identify the constructor and show the initialization call when both are available. A valid link establishes that a source range exists. Inspect whether that range also supports the answer's claim.

### Session 5 completion checks

- [ ] The answer cites the indexed commit and the correct source lines.
- [ ] An unknown evidence ID is rejected before rendering.
- [ ] An absent feature produces an evidence-scoped response.
- [ ] An answer requiring unseen startup wiring acknowledges the missing evidence.

## Session 6 — Evaluate the first phase

Create a small question set with manually labeled expected files or symbols. Run semantic search alone and then combined retrieval against the same index. Save candidate ranks separately from final answers.

| Measure | What to check |
|---|---|
| Retrieval hit rate | Whether an expected result appears in the first k candidates |
| Evidence coverage | Whether all source ranges needed for the answer were retrieved |
| Citation validity | Whether each reference points to an in-scope, verified range |
| Claim support | Whether the cited code actually establishes the statement |
| Missing-evidence handling | Whether the model names what it cannot establish |
| Cost and latency | Embedding calls, answer tokens, and end-to-end duration |

Review one failed case. If expected source was absent, improve retrieval. If it was present but the answer was wrong, inspect the prompt and evidence interpretation. Keep one successful trace and one incomplete trace using public or synthetic source.

## Phase 1 exit criteria

- [ ] A local checkout can be indexed at an explicit commit SHA.
- [ ] Uncommitted changes cannot alter that snapshot's indexed source.
- [ ] Go chunks have verified physical source ranges.
- [ ] Index publication and retries behave consistently.
- [ ] Exact and semantic search operate within the selected generation.
- [ ] LLM Playground answers a real Rate Limiter question with supporting citations.
- [ ] Skipped source and missing evidence are visible.
- [ ] Retrieval and answer quality have a small repeatable baseline.

Phase 2 adds read_file, find_symbol, and the bounded tool loop so the assistant can investigate across files after its initial retrieval.
