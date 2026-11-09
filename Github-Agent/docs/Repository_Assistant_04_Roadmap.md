# Implementation roadmap

The first two weeks deliver the Repository Assistant. The PR reviewer follows as a separate workflow and shares the source tooling. Each morning session ends with one observable deliverable.

## Week 1 · Index a commit and answer with citations

| Session | Main deliverable | Done when |
|---|---|---|
| 1 | Pinned Git snapshot CLI | Prints commit, file path, line range, and symbol for Go declarations in Rate Limiter. A dirty edit leaves output unchanged. |
| 2 | Go-aware chunks | A function, method, type, and split declaration map back to physical lines. Selected SQL and config have bounded chunks. |
| 3 | Embedding and ready generation | Same commit and profile can be indexed twice without duplicate ready rows. A simulated provider failure leaves the previous ready generation usable. |
| 4 | Exact and semantic retrieval | The limiter constructor name and natural-language query both find the expected code. Candidate ranks are logged. |
| 5 | Cited Q&A in LLM Playground | “Where is the limiter created?” returns a verified pinned link. A fabricated evidence ID is rejected. |
| 6 | Retrieval baseline | A small labeled query set compares semantic-only and hybrid retrieval. One failed answer is classified as a retrieval miss or interpretation error. |

**Week 1 demo:** Index Rate Limiter, ask one codebase question, and open the citation on the exact commit.

## Week 2 · Let the assistant investigate

| Session | Main deliverable | Done when |
|---|---|---|
| 7 | Exact `read_file` tool | Returns verified numbered source from the pinned commit. Traversal and invalid ranges are rejected. |
| 8 | `find_symbol` tool | Finds declarations and returns multiple plausible candidates when a symbol is ambiguous. |
| 9 | Playground tool integration | The model chooses a follow-up tool within fixed scope and reports when its budget is exhausted. |
| 10 | Cross-file context fixture | The agent follows handler, service, startup wiring, concrete repo, and query in both context variants. |
| 11 | Guardrail checks | Injection, unavailable source, secret filtering, and fabricated evidence are exercised. |
| 12 | Side-by-side demo and evaluation | Initial retrieval and the agent are compared on identical questions. One complete and one incomplete trace are available. |

**Week 2 demo:** Ask “How does request context reach Postgres?” on each fixture version. Verify both the actual implementation and the line where context is passed or replaced.

## Follow-on · Add PR review

| Milestone | Build | Exit criterion |
|---|---|---|
| A · PR snapshots | Add a PR URL entry point, resolve the head and merge base, and parse eligible changed lines. | The diff is bound to the reviewed snapshot and the CLI prints correct new-side anchors. |
| B · Regression reviewer | Give the reviewer scoped aliases for head and merge base. Reuse exact reads and symbol/literal search even if the head lacks embeddings. Add a finding schema and comment validation. | The reviewer flags a new dropped context, stays quiet on a clean change, and skips an unrelated PR that leaves a pre-existing dropped context in place. |
| C · Review evaluation | Compare diff-only and tool-assisted findings with seeded and clean PRs. Measure precision, recall, clean PR false positives, and invalid anchors. | The review pilot passes its quality and coverage gates in preview mode. |
| D · GitHub App | Add webhook intake, durable jobs, freshness checks, and opt-in comment publication. | One test PR receives the expected comment and replaying its delivery does not create a duplicate. |

The assistant's Q&A requires a ready semantic index for its hybrid search. Keep the shared snapshot reader and exact tools usable without a vector index so a new PR commit can still be reviewed immediately. Only index PR snapshots later if evaluation shows that semantic lookup improves findings enough to justify the additional ingestion.

The existing reviewer documents cover the separate [GitHub integration](https://chatgpt.com/library?file_id=libfile_b66946bebf3c8191899aa82ee49e9dfa&account_id=personal) and [reviewer build and operations plan](https://chatgpt.com/library?file_id=libfile_8bc6185985488191883840c5ace3dbea&account_id=personal).

After the assistant demo, prioritize follow-up work by the blocker you observe. Add GitHub URL onboarding when local setup slows adoption, directory and history tools when questions require them, and incremental indexing when full reindexing becomes expensive.

## First session checklist

- Select the local Rate Limiter checkout and resolve the intended revision to a full commit ID.
- Read tracked Go blobs from that commit, then parse functions, methods, and types.
- Print `commit · path:start_line-end_line · symbol`.
- Check three printed declarations against the committed source. Make a harmless uncommitted edit and confirm the same commit produces the same output.
