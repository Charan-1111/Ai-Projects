# Security and evaluation

## Trust boundaries and guardrails

| Risk | Control | Verification |
|---|---|---|
| Indexing uncommitted or unexpected files | Read the resolved commit's regular blobs. Allowlist source types and repo roots. | Edit and add files in the working tree. Indexed source stays the same. |
| Repository text attempts to issue instructions | Delimit snippets as untrusted data. Keep model tools read-only and validate every request in Go. | An injected source comment cannot change the tool scope, call the shell, or override citation policy. |
| Cross-repository or cross-version leakage | Bind repo and commit on the server. Include generation and access scope in all cache keys and DB queries. | Identical symbols in two repositories or commits cannot bleed into an answer. |
| Credentials or private source leave the system | Start with approved public fixtures. Apply deny paths and secret screening before embedding and before exact reads are sent to the model. Skip sensitive snippets and log a coverage reason. | A seeded token is absent from captured model requests, stored embeddings inputs, output, and logs. |
| Path traversal or resource exhaustion | Validate repository-relative paths. Bound file size, query length, tool results, calls, and deadlines. | Reject absolute paths, traversal, oversized reads, and repeated calls after exhaustion. |
| Hallucinated or misleading citations | Accept only evidence IDs issued during the current run. Assess claim support separately. | Reject unknown evidence IDs and flag valid links that fail to support the claim. |

Store model/tool timing, status, evidence IDs, and range metadata for debugging. Avoid logging raw private source or prompts. Keep the first project read-only, with app-level authorization for source calls from LLM Playground. If the snapshot store cannot serve a pinned commit, report an unavailable source range and ask the operator to restore the snapshot.

## Evaluation fixtures

| Question or scenario | Required evidence |
|---|---|
| Where is the limiter created? | Constructor and the call site that initializes it |
| Where is authentication implemented? | Middleware definition and the route registration that applies it |
| How does context reach Postgres? | Handler, service, concrete implementation wiring, and actual query argument |
| Where are connections created? | Startup pool construction and injection |
| Explain the payment flow | Route, handler, service, and provider call or queue handoff |
| Which services depend on Redis? | Client creation, dependent services, and conditional wiring |
| Context is dropped | Exact call that replaces the request context |
| Feature absent or implementation unavailable | Scoped uncertainty and the missing evidence, without a fabricated answer |

Label each fixture with an authorized commit, expected source ranges, and required answer claims. Use the Rate Limiter repository for its own behavior and a small Go fixture for other flows. Keep a clean context-propagating snapshot, one that drops context, and an unrelated change atop the broken snapshot for the later PR reviewer.

## Run two comparisons

1. **Semantic vs. hybrid retrieval.** Use the same queries and corpus. Log candidate ranks for both. Measure hit@k, the fraction of questions with at least one required result in the first k, and required-evidence coverage@k, the proportion of each question's labeled evidence returned in the first k.
2. **Initial retrieval vs. tool-assisted answers.** Fix the model version and baseline prompt policy. Record citation validity, claim support, required-claim coverage, appropriate abstentions, tool calls, tokens, and latency. Review claim support manually on the initial small fixture set.

Score path/range validity independently from whether a cited range supports the claim. Compare repeated runs if model output varies. Save one successful trace and one incomplete trace using synthetic or public code.

## Release checks

- Every emitted source link resolves to a known in-scope range from the current run.
- A dirty checkout, wrong scope, injected instruction, and seeded secret all pass the deterministic guardrail checks.
- Both context variants are distinguished correctly, including the actual concrete implementation selected at startup.
- An unavailable implementation yields a partial answer and names the unresolved step.
- The comparison identifies at least one case where a tool supplied evidence absent from the initial retrieval, or documents why the fixture did not require follow-up.
