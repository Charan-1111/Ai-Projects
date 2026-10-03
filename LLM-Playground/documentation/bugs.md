# AI Projects — Tool Calling Review

Reviewed the recent tool-related changes through commit `f5d16d4`, including PRs #44–#47.

## Overall Assessment

Your architecture now matches your original goal:

- Each service defines and executes its own tools.
- LLM-Playground discovers tools from those services.
- LLM-Playground registers their schemas with Gemini.
- When Gemini requests a tool, LLM-Playground routes the request to the appropriate service.
- Tool results are returned to Gemini to generate the final answer.

The foundation is in place, but one blocking bug and several behavior gaps need attention.

## 1. High Priority — Required Tools Keep Getting Forced

**File:** `LLM-Playground/internal/provider/gemini_provider.go`

### Problem

`geminiGenerateConfig()` sets function-calling mode to `ANY` when required tools are configured.

`Generate()` reuses the same configuration after executing the tool and returning its results.

`ANY` requires Gemini to produce a function call on every response. This can cause repeated searches until the eight-round limit or request timeout is reached.

### Current Flow

1. User asks a question.
2. Gemini is forced to call `search_documents`.
3. LLM-Playground executes the search.
4. Search results are returned to Gemini.
5. Gemini is still forced to call a tool.
6. The process repeats instead of producing a final answer.

### Required Change

Force the required tool for the initial request, then switch to `AUTO` after successful execution.

### Expected Flow

1. User asks a question.
2. Gemini calls `search_documents`.
3. LLM-Playground executes the search.
4. Search results are returned to Gemini.
5. Function-calling mode switches to `AUTO`.
6. Gemini can produce a final answer or request another tool if needed.

### Multiple Required Tools

An allowed-function-name list does not guarantee that every listed tool will execute.

If multiple distinct tools must run, track which tools have completed and explicitly enforce the remaining requirements.

Reference: https://ai.google.dev/gemini-api/docs/function-calling

## 2. High Priority — Search Parameters Are Advertised but Ignored

**Files:**

- `Semantic-Search/internal/tools/definitions.go`
- `Semantic-Search/internal/services/tool_execution.go`
- `Semantic-Search/internal/services/search.go`

### Problem

The tool schema exposes:

- `minimumScore`
- `filters.category`
- `filters.difficulty`

However, tool execution calls `SearchDocuments()`, which uses only:

- The query.
- The requested document count.

The advertised filters and similarity threshold are ignored.

### Example

A request such as:

> Find beginner Go documents with a similarity score above 0.8.

Can return documents that do not satisfy the requested category, difficulty, or threshold.

### Required Change

Choose one of these approaches:

1. Connect the tool to the chunk-search implementation that supports filtering.
2. Implement the advertised options in document search.
3. Remove unsupported parameters from the tool schema.

For your document-answering flow, retrieving relevant chunks with source metadata is the more useful direction.

Useful source metadata includes:

- Document ID.
- Document title.
- Chunk ID.
- Chunk content.
- Similarity score.

## 3. Medium Priority — Validate Tool Arguments on the Server

**Files:**

- `Semantic-Search/internal/services/tool_execution.go`
- `Semantic-Search/internal/handlers/execute_tool.go`

### Problem

The execution service unmarshals arguments but does not reject an empty query or cap `noOfDocs`.

Declaring a field as required in the model-facing schema does not validate direct HTTP requests.

### Required Validation

Before embedding or database calls:

- Reject an empty or whitespace-only query.
- Set a maximum query length.
- Apply the default document count when appropriate.
- Set a maximum document count.
- Validate the similarity threshold against your intended range.
- Validate supported filter values where applicable.

### Error Handling

The execution handler currently returns HTTP `500` for all service errors.

Use HTTP `400` for invalid arguments and unknown tool names.

Use a consistent tool-error response structure so LLM-Playground can distinguish validation failures from execution failures.

## 4. Medium Priority — Tool Registration Disables Incremental Streaming

**File:** `LLM-Playground/internal/provider/gemini_provider.go`

### Problem

`GenerateStream()` falls back to `Generate()` whenever any tool declaration exists.

It waits for the complete answer and emits a single chunk.

This affects requests even when Gemini would not use a tool.

### Required Change

For now, explicitly treat this behavior as buffered output.

Later, implement a streaming tool loop that can:

1. Receive model output.
2. Detect tool calls.
3. Execute the requested tools.
4. Return tool results to Gemini.
5. Continue streaming the final answer.

The current implementation is a reasonable intermediate step, but clients should not expect incremental tokens.

## 5. Medium Priority — Tool Discovery Failures Are Silently Skipped

**File:** `LLM-Playground/internal/tools/toolDefinitions.go`

### Problem

`RegisterTools()` logs discovery failures and continues.

The required-tool startup check helps, but optional tools can remain unavailable until the application restarts.

Duplicate tool names are also skipped. Because services are read from a map, the service that wins a duplicate-name conflict can depend on iteration order.

### Required Change

- Return a registration summary or error.
- Report unavailable services clearly.
- Reject duplicate tool names with a useful error.
- Provide a controlled retry or refresh mechanism if runtime recovery is needed.

## What Is Implemented Well

### Service-Owned Tools

Tool definitions and execution stay within their respective services.

This supports your goal of keeping LLM-Playground as the central orchestrator.

### Preserving Gemini Candidate Content

Appending the original candidate content preserves the model's tool-call context.

### Function-Call IDs

Returning function-call IDs helps match responses to the corresponding calls.

### Token Usage Across Rounds

Accumulating token usage across model requests gives a more complete view of the request's usage.

### Request Context

Forwarding the request context to tool execution supports cancellation and deadlines.

### Tool-Call Round Limit

The round limit prevents an unbounded tool-call loop.

Keep it even after fixing the forced-call configuration.

## Recommended Fix Order

1. Fix the forced-call loop.
2. Align search behavior with the tool schema.
3. Add server-side argument validation.
4. Add focused tests.
5. Improve streaming behavior.
6. Improve discovery failure handling.

## Important Test Cases

| Scenario | Expected Result |
|---|---|
| Required search succeeds | Search executes, then Gemini produces a final answer |
| Search returns no relevant results | The answer clearly reports insufficient information |
| Category and threshold are supplied | Every returned result satisfies the constraints |
| Query is blank | Request is rejected before embedding or database work |
| Document count exceeds the maximum | Request is rejected or capped according to the documented policy |
| Tool name is unknown | A clear client error is returned |
| Tool service times out | The request ends with a clear error |
| Model repeatedly requests tools | The round limit stops execution |
| Optional tool service is unavailable during startup | Its registration failure is reported clearly |
| Two services expose the same tool name | Registration reports the conflict deterministically |
| Streaming request runs with registered tools | Behavior matches the documented buffered or incremental contract |

## Review Limitations

This was a source review.

Go tests could not be run because Go was unavailable in the review environment.

Live Gemini behavior and database behavior were not verified.