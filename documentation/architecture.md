# Repository Architecture

This repository is a collection of four independently runnable Go services. The diagram shows the current local integration wiring; service URLs and external credentials are configuration-dependent.

## System overview

```mermaid
flowchart LR
    User([User or API client])

    subgraph Workspace["Ai-Projects workspace"]
        subgraph Playground["LLM-Playground · Go HTTP service · :8000"]
            LLMAPI["Generation, streaming, chat<br/>and model endpoints"]
            LLMService["Handlers and chat/generation services"]
            Provider["Gemini provider"]
            ToolRegistry["Configured HTTP tool registry"]
            MCPClient["MCP client and tool registry"]
            LLMAPI --> LLMService
            LLMService --> Provider
            LLMService --> ToolRegistry
            LLMService --> MCPClient
        end

        subgraph Semantic["Semantic-Search · Go HTTP service · configured as :8001 locally"]
            SemanticAPI["Embedding, document, chunk,<br/>search and tool endpoints"]
            SemanticService["Embedding and retrieval services"]
            Chunker["Word chunker"]
            SemanticAPI --> SemanticService
            SemanticService --> Chunker
        end

        subgraph Ask["Ask-My-Documents · Go HTTP service · :8002"]
            AskAPI["POST /rag/ask/documents"]
            AskService["RAG orchestration<br/>context + citations"]
            AskAPI --> AskService
        end

        subgraph ToolCalling["Tool-Calling · Go HTTP/MCP service · :8003"]
            ToolAPI["HTTP tool endpoints"]
            MCPServer["MCP server<br/>/mcp"]
            TimeQuestion["Gemini-backed time-question handler"]
            TimeTool["get_current_time"]
            ToolAPI --> TimeQuestion
            TimeQuestion --> TimeTool
            MCPServer --> TimeTool
        end

        SemanticDB[("PostgreSQL + pgvector<br/>documents, chunks, embeddings")]
        ChatDB[("PostgreSQL<br/>persistent conversations and messages")]
    end

    Gemini{{"Google Gemini API"}}

    User --> LLMAPI
    User --> SemanticAPI
    User --> AskAPI
    User --> ToolAPI

    ToolRegistry -->|"configured HTTP tools<br/>(e.g. document search)"| SemanticAPI
    MCPClient <-->|"MCP tool discovery and execution"| MCPServer

    AskService -->|"search request"| SemanticAPI
    SemanticAPI -->|"matching chunks"| AskService
    AskService -->|"question + retrieved context"| LLMAPI
    LLMAPI -->|"generated answer"| AskService
    AskService -->|"answer + citations"| User

    SemanticService -->|"document/chunk vectors"| SemanticDB
    SemanticDB -->|"similarity search results"| SemanticService
    LLMService -->|"conversation history"| ChatDB
    ChatDB -->|"stored messages"| LLMService

    Provider -->|"generation"| Gemini
    SemanticService -->|"embedding generation"| Gemini
    TimeQuestion <-->|"generation and function selection"| Gemini
```

The configured local integration path is LLM-Playground → Semantic-Search over HTTP and LLM-Playground → Tool-Calling over MCP. Ask-My-Documents independently coordinates Semantic-Search retrieval with LLM-Playground generation. The projects can also be run independently; these links are not a single shared process.

## Document indexing and retrieval

```mermaid
flowchart LR
    Source["Documents"] --> Upload["Semantic-Search document/chunk upload API"]
    Upload --> Split["Split content into overlapping word chunks"]
    Split --> Embed["Gemini embedding provider"]
    Embed --> Store[("PostgreSQL + pgvector")]

    Query["Search query"] --> QueryEmbed["Embed query"]
    QueryEmbed --> Similarity["Vector similarity search<br/>optional metadata filters"]
    Store --> Similarity
    Similarity --> Results["Ranked matching chunks"]
    Results --> SearchAPI["Semantic-Search response"]
```

Semantic-Search owns document and chunk indexing, embeddings, persistence, and retrieval. It exposes these capabilities through its HTTP API; it does not depend on Ask-My-Documents to perform indexing.

## Document-grounded question flow

```mermaid
sequenceDiagram
    actor User
    participant Ask as Ask-My-Documents
    participant Search as Semantic-Search
    participant LLM as LLM-Playground
    participant Gemini as Gemini API

    User->>Ask: POST /rag/ask/documents (question)
    Ask->>Search: POST /semantic/chunks/search
    Search-->>Ask: Relevant document chunks
    Ask->>Ask: Build context and source references
    Ask->>LLM: POST /v1/llm/generate (question + context)
    LLM->>Gemini: Generate grounded answer
    Gemini-->>LLM: Answer text
    LLM-->>Ask: Generated response
    Ask->>Ask: Map answer references to citations
    Ask-->>User: Answer and citations
```

## LLM tool flow

LLM-Playground can make tools available to model requests in two ways:

1. **Configured HTTP tools** are discovered and invoked through the application’s tool registry. The local configuration includes a document-search tool backed by Semantic-Search.
2. **MCP tools** are discovered from configured MCP servers and invoked through the MCP client. The local configuration connects to Tool-Calling’s `/mcp` endpoint, which exposes `get_current_time`.

Tool-Calling also has HTTP endpoints for its time and tool operations, so it can be exercised without routing through LLM-Playground.

## Service responsibilities

| Project | Owns |
| --- | --- |
| [LLM-Playground](../LLM-Playground/) | Generation, streaming, chat, Gemini provider integration, conversation persistence, and tool orchestration |
| [Semantic-Search](../Semantic-Search/) | Document ingestion, chunking, embedding, vector storage, similarity search, and search tools |
| [Ask-My-Documents](../Ask-My-Documents/) | Retrieval-augmented question orchestration, context construction, answer generation requests, and citations |
| [Tool-Calling](../Tool-Calling/) | Time-tool implementation, HTTP tool endpoints, and an MCP server exposing tools |

Each project has its own Go module, configuration, and server entrypoint. PostgreSQL connections, API credentials, and service addresses are configured per project.
