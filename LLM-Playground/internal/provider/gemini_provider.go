package provider

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"time"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client           *genai.Client
	ToolDeclarations []*genai.FunctionDeclaration
	toolExecutor     ToolExecutor
}

func NewGeminiProvider(client *genai.Client) *GeminiProvider {
	return &GeminiProvider{
		client: client,
	}
}

func (g *GeminiProvider) ConfigureTools(declarations []*genai.FunctionDeclaration, executor ToolExecutor) {
	g.ToolDeclarations = declarations
	g.toolExecutor = executor
}

func (g *GeminiProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResponse, int, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		30*time.Second,
	)
	defer cancel()

	config := geminiGenerateConfig(input, g.ToolDeclarations)
	contents := geminiContents(input)

	response, err := g.client.Models.GenerateContent(
		ctx,
		input.Model,
		contents,
		config,
	)

	if err != nil {
		classifiedErr := ClassifyError("gemini", fmt.Errorf("failed to generate response: %w", err))
		return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
	}

	var inputTokens, outputTokens, totalTokens int64
	const maxToolCallRounds = 8
	toolCalls := make([]models.ToolCall, 0)

	for round := 0; ; round++ {
		if response == nil || len(response.Candidates) == 0 {
			classifiedErr := ClassifyError("gemini", fmt.Errorf("response generation error: no candidates returned"))
			return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
		}
		if usage := response.UsageMetadata; usage != nil {
			inputTokens += int64(usage.PromptTokenCount)
			outputTokens += int64(usage.CandidatesTokenCount)
			totalTokens += int64(usage.TotalTokenCount)
		}

		candidate := response.Candidates[0]
		functionCalls := geminiFunctionCalls(candidate.Content)
		if len(functionCalls) == 0 {
			return &GenerateResponse{
				Text:         response.Text(),
				InputTokens:  inputTokens,
				OutputTokens: outputTokens,
				TotalTokens:  totalTokens,
				FinishReason: string(candidate.FinishReason),
				ToolCalls:    toolCalls,
			}, 200, nil
		}
		if round >= maxToolCallRounds {
			classifiedErr := ClassifyError("gemini", fmt.Errorf("tool call limit (%d) exceeded", maxToolCallRounds))
			return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
		}
		if g.toolExecutor == nil {
			classifiedErr := ClassifyError("gemini", fmt.Errorf("model requested a tool call but no tool executor is configured"))
			return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
		}
		if candidate.Content == nil {
			classifiedErr := ClassifyError("gemini", fmt.Errorf("model returned tool calls without candidate content"))
			return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
		}

		contents = append(contents, candidate.Content)
		functionResponses := make([]*genai.Part, 0, len(functionCalls))
		for _, functionCall := range functionCalls {
			toolExecutionStart := time.Now()

			toolResult, executionErr := g.toolExecutor.ToolExecution(ctx, functionCall.Name, functionCall.Args)
			toolExecutionDuration := time.Since(toolExecutionStart).Milliseconds()

			if executionErr != nil {
				classifiedErr := ClassifyError("tool", fmt.Errorf("execute tool %q: %w", functionCall.Name, executionErr))
				return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
			}
			functionResponses = append(functionResponses, &genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					ID:       functionCall.ID,
					Name:     functionCall.Name,
					Response: toolResult,
				},
			})
			toolCalls = append(toolCalls, models.ToolCall{
				Name:      functionCall.Name,
				Arguments: functionCall.Args,
				Duration:  int32(toolExecutionDuration),
				Status:    "success",
			})
		}
		contents = append(contents, &genai.Content{Role: "user", Parts: functionResponses})

		response, err = g.client.Models.GenerateContent(ctx, input.Model, contents, config)
		if err != nil {
			classifiedErr := ClassifyError("gemini", fmt.Errorf("failed to generate response after tool execution: %w", err))
			return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
		}
	}
}

func (g *GeminiProvider) GenerateStream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, <-chan error) {
	if len(g.ToolDeclarations) > 0 {
		chunks := make(chan StreamChunk)
		errs := make(chan error, 1)
		go func() {
			defer close(chunks)
			defer close(errs)
			response, _, err := g.Generate(ctx, input)
			if err != nil {
				errs <- err
				return
			}
			select {
			case chunks <- StreamChunk{
				Delta:        response.Text,
				InputTokens:  response.InputTokens,
				OutputTokens: response.OutputTokens,
				FinishReason: response.FinishReason,
				ToolCalls:    response.ToolCalls,
			}:
			case <-ctx.Done():
			}
		}()
		return chunks, errs
	}

	ctx, cancel := context.WithTimeout(
		ctx,
		30*time.Second,
	)

	chunks := make(chan StreamChunk)
	errs := make(chan error, 1)

	config := geminiGenerateConfig(input, g.ToolDeclarations)
	contents := geminiContents(input)

	go func() {
		defer cancel()
		defer close(chunks)
		defer close(errs)
		defer func() {
			if r := recover(); r != nil {
				errs <- ClassifyError("gemini", fmt.Errorf("gemini stream: recovered from panic: %v", r))
			}
		}()

		streamChunks := g.client.Models.GenerateContentStream(
			ctx,
			input.Model,
			contents,
			config,
		)

		for chunk, err := range streamChunks {
			if err != nil {
				errs <- ClassifyError("gemini", fmt.Errorf("gemini stream error: %w", err))
				return
			}

			if chunk == nil || len(chunk.Candidates) == 0 {
				continue
			}

			streamChunk := StreamChunk{
				Delta:        chunk.Text(),
				FinishReason: string(chunk.Candidates[0].FinishReason),
			}
			if chunk.UsageMetadata != nil {
				streamChunk.InputTokens = int64(chunk.UsageMetadata.PromptTokenCount)
				streamChunk.OutputTokens = int64(chunk.UsageMetadata.CandidatesTokenCount)
			}

			select {
			case chunks <- streamChunk:
			case <-ctx.Done():
				return
			}
		}
	}()

	return chunks, errs
}

func geminiFunctionCalls(content *genai.Content) []*genai.FunctionCall {
	if content == nil {
		return nil
	}

	functionCalls := make([]*genai.FunctionCall, 0)
	for _, part := range content.Parts {
		if part != nil && part.FunctionCall != nil {
			functionCalls = append(functionCalls, part.FunctionCall)
		}
	}
	return functionCalls
}

func geminiContents(input GenerateInput) []*genai.Content {
	contents := make([]*genai.Content, 0, len(input.History)+1)
	for _, message := range input.History {
		if message.Content == "" {
			continue
		}
		role := message.Role
		if role == "assistant" {
			role = "model"
		}
		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: message.Content}},
		})
	}

	if input.Prompt != "" {
		contents = append(contents, genai.Text(input.Prompt)...)
	}

	return contents
}

func geminiGenerateConfig(input GenerateInput, declarations []*genai.FunctionDeclaration) *genai.GenerateContentConfig {
	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(input.Temperature)),
	}

	if input.SystemPrompt != "" {
		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: input.SystemPrompt}},
		}
	}
	if len(declarations) > 0 {
		config.Tools = []*genai.Tool{{FunctionDeclarations: declarations}}
	}
	if input.MaxOutputTokens > 0 {
		config.MaxOutputTokens = int32(input.MaxOutputTokens)
	}
	if len(input.RequiredTools) > 0 {
		config.ToolConfig = &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode:                 genai.FunctionCallingConfigModeAny,
				AllowedFunctionNames: append([]string(nil), input.RequiredTools...),
			},
		}
	}

	return config
}
