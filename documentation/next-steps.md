# Next Learning Milestone — Multi-Tool Assistant

## Goal

Build a small multi-tool assistant using your existing services.

This will help you learn how an LLM:

- Chooses the appropriate tool.
- Uses multiple tools for one request.
- Combines tool results into one answer.
- Uses the output of one tool as input to another.

Keep fixing the existing bugs in parallel. You do not need to start a separate project yet.

## Step 1 — Combine Your Existing Tools

Expose both tools through LLM-Playground:

| Tool | Owning Service | Purpose |
|---|---|---|
| `get_current_time` | Tool-Calling | Get the current time in a timezone |
| `search_documents` | Semantic-Search | Retrieve relevant documents |

### Test Requests

| Request | Expected Behavior |
|---|---|
| “What time is it in India?” | Calls `get_current_time` |
| “How do Go channels work?” | Calls `search_documents` |
| “What time is it in India, and find documents explaining Go channels?” | Calls both tools and combines their results |

### Configuration

Use a configuration where tools are optional for this experiment.

Your retrieval-only configuration can continue requiring search once the forced-call bug is fixed.

### Deliverable

One user request can use multiple services and produce one combined answer.

## Step 2 — Add a Tool-Execution Trace

Return a small execution trace alongside the answer.

### Example Response

{
  "answer": "The current time is ... Go channels allow ...",
  "tool_calls": [
    {
      "name": "get_current_time",
      "arguments": {
        "timezone": "Asia/Kolkata"
      },
      "duration_ms": 35,
      "status": "success"
    },
    {
      "name": "search_documents",
      "arguments": {
        "query": "Go channels"
      },
      "duration_ms": 420,
      "status": "success"
    }
  ]
}

### What This Helps You Understand

- Which tools the model selected.
- Which arguments it generated.
- Which tools your application executed.
- How long each execution took.
- Whether each execution succeeded.

### Deliverable

You can inspect a request without relying on scattered console logs.

## Step 3 — Add a Dependent Tool Sequence

Introduce a tool whose input comes from another tool's result.

### Tools

| Tool | Input | Output |
|---|---|---|
| `search_documents` | Search query | Document IDs and relevant excerpts |
| `get_document` | Document ID | Full document |

### Test Request

> Find the document about Go cancellation, then explain its complete example.

### Expected Flow

1. Gemini calls `search_documents`.
2. Semantic-Search returns matching documents with document IDs.
3. Gemini selects a document ID from those results.
4. Gemini calls `get_document` with that ID.
5. The service returns the full document.
6. Gemini explains the example using the retrieved content.

Ensure that returned document IDs correspond to real documents that `get_document` can fetch.

### Deliverable

A successful dependent sequence:

Search → Fetch document → Answer

## Step 4 — Introduce MCP

After the multi-tool and dependent-tool flows work:

1. Expose one existing service through MCP.
2. Connect it to LLM-Playground.
3. Compare it with your existing custom HTTP tool protocol.

### Compare These Areas

- Tool discovery.
- Tool schemas.
- Tool execution.
- Error representation.

### Deliverable

One existing service is accessible through MCP, and you understand how its integration differs from your custom protocol.

## Next Morning Session

Focus only on Step 1.

### Main Task

Make `get_current_time` and `search_documents` work together through LLM-Playground.

### Completion Criteria

- A time-only request uses the time tool.
- A document question uses the search tool.
- A combined request uses both tools.
- The user receives one combined answer.

## Suggested Order

1. Combine the existing tools.
2. Add execution traces.
3. Add a dependent tool sequence.
4. Introduce MCP.

Continue fixing the existing issues alongside these milestones.