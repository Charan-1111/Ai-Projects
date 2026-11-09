# GitHub Repository Assistant
## Product requirements

**Status** Proposed · **Initial repository** Rate Limiter · **Delivery** Twelve morning sessions over two weeks

Build a Go service inside AI Projects that answers questions about a specific commit and cites the exact source lines. LLM Playground remains the user-facing model and tool orchestrator. The Repository Assistant owns indexing and scoped source tools. A later Code Review Agent uses those tools to investigate a pull request across snapshots.

| Document | Use it for |
|---|---|
| [02 · Technical design](Repository_Assistant_02_Technical_Design.md) | Service boundaries, ingestion, storage, ranking, tool contracts, and citation validation |
| [03 · Security and evaluation](Repository_Assistant_03_Security_and_Evaluation.md) | Guardrails, test fixtures, and measures of retrieval and answer quality |
| [04 · Roadmap](Repository_Assistant_04_Roadmap.md) | Twelve build sessions, release gates, and the PR review extension |

## User flow

1. Select a trusted local checkout and resolve a revision to a commit.
2. Index the tracked source files at that commit. Report skipped files and ingestion status.
3. Ask a question in LLM Playground with a selected repository and snapshot.
4. Receive an answer with verified references or a clear explanation of which evidence is missing.

The first demo uses the Rate Limiter repository for real questions. A small Go fixture supplies routing, database wiring, Redis, and payment examples that the Rate Limiter repository may lack.

## Functional requirements

| ID | Requirement | Acceptance criterion |
|---|---|---|
| R1 | Index an immutable Git snapshot | An uncommitted change cannot alter the content indexed under a commit ID. The selected commit appears in output. |
| R2 | Extract Go declarations and relevant source | Functions, methods, and types have correct physical line ranges. Large declarations split without changing the cited ranges. Selected SQL and config remain searchable. |
| R3 | Publish complete index generations | A repeated index of the same inputs creates no duplicate ready generation. A failed reindex leaves the earlier ready generation available. |
| R4 | Support exact and semantic discovery | Both a constructor name and a natural-language limiter question retrieve expected evidence. |
| R5 | Answer with source references | Each factual code claim references evidence from the current run. Unknown evidence IDs cannot produce links. |
| R6 | Follow relationships across files | Given a handler, interface, and concrete repository, the assistant can inspect wiring and trace the actual query call. |
| R7 | Explain evidence gaps | A missing implementation, ambiguous binding, or exhausted budget produces a partial answer with the unresolved step named. |
| R8 | Scope reads and model input | Tools cannot switch repositories or snapshots. Denied source is excluded from embeddings, tool output, and prompts. |

## Scope boundaries

The MVP starts with local Git ingestion, Go source plus selected SQL/config, PostgreSQL with pgvector, and public or explicitly authorized source. The first two weeks cover ingestion, cited Q&A, and a bounded read-only investigation loop.

Remote onboarding, incremental reindexing, directory exploration, history tools, and pull request publishing are later milestones. Start without a full call graph, type-aware resolution, shell execution, or a dedicated frontend. Record the actual search scope when answering questions about repository-wide absence.

## Completion criteria

By the end of Week 1, “Where is the limiter created?” returns a commit-pinned citation from the real repository. By the end of Week 2, the assistant follows handler → service → concrete repository → query in the fixture and distinguishes the version that passes the request context from the version that replaces it. The evaluation report compares the initial retrieval answer with the tool-assisted answer and preserves evidence traces for one success and one incomplete investigation.
