# AI / LLM Project-Based Learning Roadmap

## Should I Stop the LLM Playground Here?

Yes. **Stop adding major features to the LLM Playground for now and start Project 2.**

The reason is not that the Playground is fully production-complete. The reason is that it has already taught the major concepts it was designed to teach.

Since the goal is to learn AI and LLM concepts through projects, the objective should not be:

> Make every project production-complete.

Instead:

> Build enough to understand the concepts, then move to a project that forces you to learn the next set of concepts.

---

# Project 1 — LLM Playground ✅

## Concepts Covered

The LLM Playground has already covered a large part of the basic LLM engineering layer:

- What an LLM API looks like
- Gemini API integration
- Non-streaming generation
- Temperature
- Maximum output tokens
- Token usage
- Cost estimation
- Latency measurement
- Streaming
- Server-Sent Events (SSE)
- Context cancellation
- Timeouts
- Retries
- Exponential backoff
- Jitter
- Provider abstraction
- Provider error classification
- Rate-limit handling
- System instructions
- Conversation history
- In-memory conversations
- Persistent conversations
- Database-backed conversation history

The learning progression has roughly been:

```text
LLM
 ↓
Tokens / Context
 ↓
LLM API
 ↓
Temperature / Max Tokens
 ↓
Streaming
 ↓
Usage / Cost / Latency
 ↓
Timeout / Cancellation
 ↓
Retries / Backoff / Jitter
 ↓
Provider Abstraction
 ↓
System Instructions
 ↓
Conversation History
 ↓
Persistence

        YOU ARE HERE
              │
              ▼
       Start Project 2
```

There are still features that could be added, such as structured outputs, conversation summarization, additional providers, tests, and deeper observability.

However, continuing to add features will increasingly teach **backend hardening** rather than fundamentally new AI concepts.

Those improvements can be revisited later.

---

# Before Leaving Project 1

Spend one final session cleaning up the project.

## README

Document:

- What the project does
- Architecture
- API endpoints
- Gemini integration
- Streaming architecture
- Retry strategy
- Conversation architecture
- Database architecture
- Concepts learned
- How to run the project

## Future Improvements

Keep unfinished ideas under a section such as:

```text
Future Improvements

- Structured outputs
- Context summarization
- Token-budget-based conversation management
- Multiple LLM providers
- Better observability
- Unit tests
- Integration tests
- Model routing
```

Then freeze the project and move forward.

---

# Project 2 — Developer Knowledge Assistant

This should be the next major project.

The goal is to learn:

- Embeddings
- Vector mathematics
- Semantic similarity
- Cosine similarity
- Vector databases
- pgvector
- Chunking
- Semantic search
- Top-K retrieval
- RAG
- Reranking
- RAG evaluation

Do **not** start by building the entire RAG system.

Build it incrementally.

---

# Phase 1 — Learn Embeddings

Start with the smallest possible system:

```text
Text
 ↓
Embedding Model
 ↓
Vector
```

Example:

```text
"Go is a programming language"

        ↓

Embedding Model

        ↓

[0.12, -0.38, 0.91, ...]
```

## Concepts to Learn

Understand:

- What is an embedding?
- Why do we represent text as vectors?
- What is an embedding dimension?
- What is semantic similarity?
- What is cosine similarity?
- What is dot product?
- What does vector normalization mean?
- How does an embedding model differ from a generative LLM?

## First Implementation

Build something like:

```text
POST /embeddings
```

Request:

```json
{
  "text": "Goroutines are lightweight threads"
}
```

Flow:

```text
Go API
 ↓
Embedding API
 ↓
Embedding Vector
```

Response conceptually:

```json
{
  "dimensions": 768,
  "embedding": [0.12, -0.42, 0.81]
}
```

The exact number of dimensions depends on the embedding model being used.

---

# Phase 2 — Compare Embeddings

Now take multiple sentences.

For example:

```text
A = "How do I reset my password?"

B = "I forgot my login credentials."

C = "How do I cook biryani?"
```

Generate embeddings:

```text
A → Vector A
B → Vector B
C → Vector C
```

Then calculate similarity.

Expected conceptual result:

```text
similarity(A, B) → HIGH

similarity(A, C) → LOW
```

Implement cosine similarity yourself in Go before relying entirely on a vector database.

That will help you understand what the database is eventually doing for you.

---

# Phase 3 — Semantic Search

Once embeddings make sense, introduce document storage.

```text
Documents
    ↓
Chunking
    ↓
Embeddings
    ↓
PostgreSQL + pgvector
    ↓
Store Vectors
```

Then implement search:

```text
User Question
      ↓
Create Embedding
      ↓
Vector Similarity Search
      ↓
Top-K Chunks
```

## Concepts to Learn

- Vector databases
- pgvector
- Similarity search
- Top-K
- Chunk size
- Chunk overlap
- Metadata
- Metadata filtering
- Vector indexes
- Exact vs approximate search

At this point you have built **semantic search**, but not RAG yet.

---

# Phase 4 — Build RAG

Once semantic search works, connect the LLM.

```text
Question
   ↓
Embedding
   ↓
Vector Search
   ↓
Top-K Relevant Chunks
   ↓
Build Context
   ↓
Gemini
   ↓
Grounded Answer
```

RAG becomes easier to understand when you realize it is fundamentally:

> Retrieve relevant information → put that information into the LLM context → ask the LLM to answer using it.

## Developer Knowledge Assistant Architecture

```text
                  Documents
                      │
         ┌────────────┼────────────┐
         ↓            ↓            ↓
      Go Docs    Backend Notes   PDFs / Markdown
         │            │            │
         └────────────┼────────────┘
                      ↓
                   Chunker
                      ↓
                Embedding Model
                      ↓
              PostgreSQL + pgvector


User:
"How does context cancellation work in Go?"

                      ↓
                Query Embedding
                      ↓
                Vector Search
                      ↓
              Relevant Chunks
                      ↓
                    Gemini
                      ↓
              Answer + Sources
```

---

# Phase 5 — Improve RAG

Once basic RAG works, learn why naive RAG can fail.

Add concepts gradually:

- Better chunking strategies
- Metadata filtering
- Query rewriting
- Hybrid search
- Reranking
- Context construction
- Citations
- Retrieval evaluation
- Answer evaluation

Do not add all of these at once.

Use the same process:

```text
Learn
 ↓
Implement
 ↓
Observe
 ↓
Break
 ↓
Understand
 ↓
Improve
```

---

# Project 3 — Tool-Calling Developer Assistant

After RAG, move to tool calling.

The goal is to understand how an LLM interacts with external systems.

Example:

```text
User
 ↓
"Show me issue #42 from my repository"

LLM
 ↓
Decides it needs a tool

 ↓

get_github_issue({
    "issue_number": 42
})

 ↓

Go Backend Executes Tool

 ↓

GitHub API

 ↓

Tool Result

 ↓

LLM

 ↓

Final Answer
```

## Concepts to Learn

- Structured outputs
- Function calling
- Tool schemas
- Tool arguments
- Tool execution
- Returning tool results
- Multiple tools
- Parallel tool calls
- Tool errors
- Tool authorization

Possible tools:

```text
get_github_issue()
search_documentation()
query_database()
get_pull_request()
search_repository()
```

---

# Project 4 — AI Agent

Only after tool calling feels comfortable should you move to agents.

An agent can conceptually be thought of as:

```text
while goal_not_complete {

    observe()

    reason()

    choose_action()

    execute_tool()

    observe_result()
}
```

Architecture:

```text
User Goal
   ↓
LLM
   ↓
Decide Next Action
   ↓
Tool
   ↓
Observation
   ↓
LLM
   ↓
Decide Again
   ↓
...
   ↓
Final Answer
```

## Concepts to Learn

- Agent loop
- Agent state
- Short-term memory
- Long-term memory
- Planning
- ReAct
- Routing
- Reflection
- Human-in-the-loop
- Checkpoints
- MCP
- Agent evaluation
- Multi-agent systems

Do not immediately use an agent framework.

First implement a simple loop yourself so you understand what the framework is abstracting.

---

# Project 5 — Production AI Platform

Once LLMs, RAG, tools, and agents make sense, focus on production AI engineering.

## Concepts

- LLM observability
- Tracing
- Evaluation
- Guardrails
- Prompt injection
- Tool security
- Data leakage prevention
- Caching
- Model routing
- Model fallbacks
- Cost optimization
- Token budgets
- Latency optimization
- Rate limiting
- Deployment
- Monitoring

This is where backend engineering and AI engineering strongly come together.

---

# Complete Project Progression

```text
PROJECT 1 — LLM Playground                         ✅
────────────────────────────────────────────────────
LLM fundamentals
Tokens / context
LLM APIs
Generation
Streaming
Prompting
Reliability
Conversation
Persistence

                    ↓

PROJECT 2 — Developer Knowledge Assistant           ← NEXT
────────────────────────────────────────────────────
Embeddings
Vector mathematics
Cosine similarity
Vector databases
pgvector
Chunking
Semantic search
Top-K
RAG
Reranking
RAG evaluation

                    ↓

PROJECT 3 — Tool-Calling Developer Assistant
────────────────────────────────────────────────────
Structured outputs
Function calling
Tool schemas
Tool execution
Multiple tools
Parallel tools
Tool errors

                    ↓

PROJECT 4 — AI Agent
────────────────────────────────────────────────────
Agent loop
State
Memory
Planning
ReAct
Routing
Reflection
Human-in-the-loop
Checkpoints
MCP
Agent evaluation
Multi-agent patterns

                    ↓

PROJECT 5 — Production AI Platform
────────────────────────────────────────────────────
Observability
Tracing
Evals
Guardrails
Prompt injection
Security
Caching
Model routing
Fallbacks
Cost optimization
Deployment
```

---

# When Should I Stop a Project?

Do not ask:

> Is this project completely finished?

Instead ask:

> Have I understood the concepts this project was designed to teach?

If the answer is yes:

```text
Commit the project
      ↓
Clean the README
      ↓
Document the architecture
      ↓
Document what you learned
      ↓
Write future improvements
      ↓
Tag / version it
      ↓
MOVE ON
```

You can always return later.

---

# Recommended Learning Pattern

For every concept:

```text
1. Understand WHY the concept exists

                ↓

2. Learn the basic theory

                ↓

3. Implement the smallest version

                ↓

4. Test it

                ↓

5. Break it intentionally

                ↓

6. Understand the failure

                ↓

7. Improve the implementation

                ↓

8. Move to the next concept
```

Avoid this pattern:

```text
Watch 20 hours of AI courses
        ↓
Take lots of notes
        ↓
Learn 30 concepts
        ↓
Try to build something
```

Prefer:

```text
Concept
 ↓
Code
 ↓
Experiment
 ↓
Understand
 ↓
Next Concept
```

---

# Immediate Next Step

Freeze the LLM Playground after a small cleanup session.

Then start the Developer Knowledge Assistant with this first milestone:

```text
DAY 1

Learn:
- What is an embedding?
- Why vectors?
- Embedding dimensions
- Semantic similarity
- Cosine similarity

Build:

POST /embeddings

        ↓

Embedding API

        ↓

Vector
```

Then:

```text
DAY 2

Generate embeddings for multiple sentences

        ↓

Implement cosine similarity in Go

        ↓

Compare semantic similarity
```

Then:

```text
DAY 3+

Documents
 ↓
Chunking
 ↓
Embeddings
 ↓
pgvector
 ↓
Semantic Search
 ↓
RAG
```

The overall learning philosophy should remain:

> **Concept → Implement → Observe → Understand → Improve → Move On**

This keeps the projects focused on learning AI concepts while continuously using and strengthening Go/backend engineering skills.
