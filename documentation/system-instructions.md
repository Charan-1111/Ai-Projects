# System Instructions in the LLM Playground

## What Is a System Instruction?

A system instruction tells the model how it should behave before it receives the user's prompt.

Example:

```text
You are a helpful Go programming assistant.
Always provide idiomatic Go examples.
Keep explanations suitable for beginners.
```

The user prompt could then be:

```text
Explain interfaces in Go.
```

The model uses both inputs:

- **System instruction:** defines the assistant's behavior.
- **User prompt:** contains the user's actual question.

## Why Are System Instructions Useful?

### Personality and tone

```text
You are a concise technical assistant.
Avoid unnecessary explanations.
```

### Audience

```text
Explain all concepts to someone who is new to programming.
```

### Output format

```text
Return the answer as JSON with the fields:
summary, example, and caveats.
```

### Domain behavior

```text
You are an AI assistant specialized in backend Go development.
Prefer standard library solutions before suggesting external packages.
```

### Safety and boundaries

```text
Do not expose API keys, passwords, or internal configuration.
```

System instructions make these rules reusable. Without them, every client must repeat the same rules in every user prompt.

## Current Project State

The project already contains part of the required design. In `LLM-Playground/internal/provider/providerModel.go`, `GenerateInput` has a `SystemPrompt` field:

```go
type GenerateInput struct {
    SystemPrompt    string
    Prompt          string
    Model           string
    Temperature     float64
    MaxOutputTokens int64
}
```

However, the field is not currently populated or sent to Gemini.

`BuildGenerateInput` currently copies the user prompt:

```go
input.Prompt = request.Prompt
```

The Gemini provider sends only the user prompt:

```go
genai.Text(input.Prompt)
```

Therefore, system instructions currently have no effect.

## Is It Necessary for This Project?

System instructions are not required for the basic endpoint to work. The project can already:

- receive a prompt,
- call Gemini,
- stream responses,
- retry provider failures,
- calculate token usage and cost.

They become useful as the playground grows into a reusable LLM service. They can provide:

- consistent assistant behavior,
- specialized model roles,
- predictable response formats,
- easier testing of different model behaviors,
- common instructions shared by all clients.

For example, different model configurations could represent different roles:

```json
{
  "available_models": {
    "go-tutor": {
      "provider_model": "gemini-3.5-flash-lite",
      "system_instruction": "You are a beginner-friendly Go programming tutor."
    },
    "code-reviewer": {
      "provider_model": "gemini-3.7-flash",
      "system_instruction": "Review code for correctness, security, and maintainability."
    }
  }
}
```

## Recommended Integration

### 1. Add an instruction to the model configuration

Update `LLM-Playground/internal/models/model_config.go`:

```go
type ModelConfig struct {
    DisplayName       string `json:"display_name"`
    ProviderModel     string `json:"provider_model"`
    SystemInstruction string `json:"system_instruction"`

    InputCostPerMillionTokens  float64 `json:"input_cost_per_million_tokens"`
    OutputCostPerMillionTokens float64 `json:"output_cost_per_million_tokens"`
    MaxOutputTokens            int64   `json:"max_output_tokens"`
    SupportsStreaming          bool    `json:"supports_streaming"`
}
```

### 2. Copy the instruction into the provider input

Update `BuildGenerateInput`:

```go
input.SystemPrompt = modelConfig.SystemInstruction
```

This creates the flow:

```text
model_id -> model configuration -> system instruction + provider model
```

### 3. Pass the instruction to Gemini

Gemini supports system instructions through `GenerateContentConfig`. The provider should add the instruction to the configuration used by both normal and streaming generation:

```go
config := &genai.GenerateContentConfig{
    Temperature: genai.Ptr(float32(input.Temperature)),
}

if input.SystemPrompt != "" {
    config.SystemInstruction = genai.NewContentFromText(
        input.SystemPrompt,
        genai.RoleUser,
    )
}
```

The exact SDK helper should be checked against the version in `LLM-Playground/go.mod` before implementation.

## Configuration-Based or Request-Based?

There are two possible designs.

### Configuration-based instruction

```json
{
  "prompt": "Explain Go interfaces",
  "model_id": "go-tutor"
}
```

The server gets the instruction from the selected model configuration.

This is better for:

- consistent behavior,
- trusted application rules,
- production use,
- preventing clients from replacing important instructions.

### Request-based instruction

```json
{
  "system_prompt": "Answer like a Go teacher",
  "prompt": "Explain interfaces"
}
```

This is useful for an experimentation playground because users can test different instructions. However, the behavior can change on every request, so it should not be trusted for security or authorization rules.

## Recommendation

Use model-configuration instructions first. The project already selects models through `model_id`, so this approach fits the existing architecture.

An optional request-level instruction can be added later for experimentation, but it should not replace important application rules.

System instructions are also not a complete security boundary. A user can ask the model to ignore them, so authentication, authorization, validation, and protection of secrets must remain in Go code.

## Summary

System instructions are not necessary for the current playground to function. They are the natural next feature for making model behavior consistent and configurable.

The project has already started the integration by defining `GenerateInput.SystemPrompt`. The remaining work is to populate that field from model configuration and pass it to Gemini's generation configuration for both normal and streaming requests.
