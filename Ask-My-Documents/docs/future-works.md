# Ask My Documents — Review and Next Steps

Reviewed against `main` at commit [`3226ca2`](https://github.com/Charan-1111/Ai-Projects/commit/3226ca2fcc5834f18f0d4233492d2326b5cba382).

## Current progress

Your service now has the main pieces of a basic RAG pipeline:

1. Accept a question.
2. Call Semantic Search for relevant chunks.
3. Build numbered excerpts.
4. Send the question and excerpts to LLM Playground.
5. Return an answer and citations.

The excerpt numbering is consistent, and citations are now limited to sources whose reference numbers appear in the generated answer.

## Fix first: the LLM response has an outer wrapper

The `/v1/llm/generate` endpoint returns an object shaped like this:

```json
{
  "code": 0,
  "message": "Response generated successfully",
  "response": {
    "request_id": "req-123",
    "model": "gemini-3.5-flash-lite",
    "response": "According to your notes...",
    "usage": {
      "input_tokens": 100,
      "output_tokens": 30,
      "total_tokens": 130
    }
  }
}
```

Your current client unmarshals the **entire body** directly into `LLMResponse`. However, `LLMResponse` describes only the **inner** `response` object. Consequently, its answer field remains empty even when generation succeeds.

Add an envelope type:

```go
type GenerateAPIResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Response LLMResponse `json:"response"`
}
```

Then decode the envelope in `LLMClient.Generate`:

```go
var result models.GenerateAPIResponse

if err := sonic.Unmarshal(resp, &result); err != nil {
	return nil, fmt.Errorf("decode LLM response: %w", err)
}

if result.Code != 0 {
	return nil, fmt.Errorf("LLM API error: %s", result.Message)
}

return &result.Response, nil
```

Your service can then read `llmResponse.Response` as it currently does.

> After making this change, send one real request through all three services: Ask My Documents → Semantic Search → LLM Playground. Verify that the returned `answer` is nonempty.

## Remaining correctness issues

### 1. Validate every reference the LLM writes

Currently, you scan known source numbers and include a source if its `[n]` appears in the answer. That handles valid references but misses invented ones.

For example:

```text
Context supplied: [1], [2], [3]
LLM answer: "Your notes use a token bucket [9]."
```

The answer still contains `[9]`, but the `citations` array cannot resolve it.

**Next step:** Extract all `[number]` references from the answer. For each one:

1. Look it up in `sources`.
2. Add it to the returned citations if found.
3. If it does not exist, reject the generated answer or retry generation with an instruction to use only supplied references.

Also remember that a valid reference number does **not** prove the cited chunk supports the claim. Check that separately during evaluation.

### 2. Fix partial search filters

Ask My Documents accepts `category` and `difficulty` independently. Semantic Search currently chooses a filtered SQL query if **either** field is supplied, but that query requires **both** fields to match.

That means a category-only request can return no results even when matching documents exist.

Make each condition optional in the SQL query, or support separate query paths for:

- No filters
- Category only
- Difficulty only
- Both filters

Test all four cases.

### 3. Fix the server error channel

The server receives once from `listenErr` in a `select`, then tries to receive from it again while logging. The second receive can block forever.

Use the error received by the `select`:

```go
case err := <-listenErr:
	app.log.Log.Error().Err(err).Msg("server failed to listen")
```

### 4. Stop configuration loading after a read failure

If `os.ReadFile` fails, `LoadConfig` currently continues to unmarshal the returned bytes. Return from the callback immediately:

```go
fileBytes, err := os.ReadFile(configPath)
if err != nil {
	loadErr = err
	return
}

loadErr = sonic.Unmarshal(fileBytes, c)
```

### 5. Tighten input and context limits

- Reject prompts containing only whitespace.
- Set a maximum for `noOfDocs`.
- Set a maximum size for the combined excerpts sent to the LLM.
- Decide what `minimumScore: 0` means. It currently becomes the default `0.7`, so callers cannot explicitly request zero. Use a pointer field if you need to distinguish “omitted” from “set to zero.”

## Suggested tests

Start with a small set of integration questions against notes whose contents you know:

| Test | Expected behavior |
| --- | --- |
| Answer appears in one chunk | Answer cites that chunk |
| Answer needs two chunks | Both supporting chunks are cited |
| Answer is absent from all documents | Says it could not find the information |
| LLM returns `[9]`, but only `[1]`–`[3]` exist | Invented citation is caught |
| Category-only filter | Finds matching documents |
| Semantic Search is unavailable | Returns an error, not an empty success |
| LLM response is valid JSON with no answer | Returns an error |

For a wrong answer, inspect the stages separately: **Were the right chunks retrieved? Did the prompt include them? Did the LLM use them correctly?**

## Future improvements

Build these after the basic flow works reliably:

### Retrieval quality

- **Hybrid search:** Combine vector similarity with keyword search. This helps with exact identifiers, function names, and error messages that embeddings may overlook.
- **Reranking:** Retrieve more candidates, then reorder them by relevance before selecting excerpts for the prompt.
- **Adjacent chunks:** When an answer spans a chunk boundary, include a neighboring chunk from the same document.
- **Document filters:** Let users restrict a question to selected documents or projects.

### Answer quality

- **Context budget:** Select excerpts based on a token budget rather than a fixed number of chunks.
- **Clearer source locations:** Save section headings, page numbers, or file paths during ingestion and include them in citations.
- **Grounding evaluation:** Check whether each answer claim is supported by its cited excerpt.
- **Prompt injection tests:** Upload a document containing instructions such as “ignore the question” and confirm document content does not control the assistant.

### Project reliability

- **Observability:** Record retrieval time, LLM time, selected chunk IDs, token usage, and errors. Avoid logging full private document contents by default.
- **Ingestion status:** Expose whether a document is pending, indexed, or failed.
- **Reindexing:** Support replacing a document and its old chunks safely.
- **Evaluation set:** Keep a small set of questions with expected answers and source documents. Run it whenever you change chunking, retrieval, or prompting.

## Recommended order

1. Fix the `/generate` response envelope.
2. Test one complete question-to-answer request.
3. Validate cited reference numbers.
4. Fix partial search filters and the server/configuration bugs.
5. Add a small RAG evaluation set.
6. Experiment with context budgets, hybrid search, and reranking.

**Milestone for the first working version:** A question about your notes returns a supported answer with correct chunk citations; a question absent from the notes returns an honest “not found” response.