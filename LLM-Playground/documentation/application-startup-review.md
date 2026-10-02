# Application Startup Review

## Current Assessment

The overall startup wiring is reasonable. The chat services and `Application.provider` share the same `*GeminiProvider`, so assigning tool declarations in `StartServer` updates the provider used for generation.

## Findings

### Tool discovery failures do not fail startup

`StartServer` calls `RegisterTools` and continues to start the server. `RegisterTools` logs fetch or decode errors and continues, so Gemini may start with an empty or incomplete declaration list.

**Recommendation:** Have `RegisterTools` return an error (or a structured result) and decide explicitly whether tool discovery failures should prevent startup.

### Gemini provider assertion is brittle

`NewApplication` always constructs a Gemini provider, but `StartServer` checks that assumption with a runtime type assertion.

**Recommendation:** Keep a typed `*provider.GeminiProvider` during construction, or add a provider method for setting declarations. Configure it before injecting the provider into chat services where practical.

### Some `Application` fields appear redundant

`client`, `clients`, and `apiFactory` are stored on `Application` but are not referenced outside construction. The client is needed to create the provider; the HTTP clients and factory are needed to create the tool registry.

**Recommendation:** Remove these fields if no later code needs them; keep the local variables used during construction.

## Function-Calling Follow-Up

Passing declarations allows Gemini to choose a function, but does not execute the returned call. Complete the flow by detecting Gemini's function-call part, invoking the matching registered tool, and sending a `FunctionResponse` back to Gemini so it can produce the final answer.