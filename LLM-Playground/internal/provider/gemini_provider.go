package provider

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	Client *genai.Client
}

func (g *GeminiProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResponse, int, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		30*time.Second,
	)
	defer cancel()

	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(input.Temperature)),
	}

	if input.MaxOutputTokens > 0 {
		config.MaxOutputTokens = int32(input.MaxOutputTokens)
	}

	response, err := g.Client.Models.GenerateContent(
		ctx,
		input.Model,
		genai.Text(input.Prompt),
		config,
	)

	if err != nil {
		classifiedErr := ClassifyError("gemini", fmt.Errorf("failed to generate response: %w", err))
		return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
	}

	if response == nil || response.UsageMetadata == nil {
		classifiedErr := ClassifyError("gemini", fmt.Errorf("response generation error"))
		return &GenerateResponse{}, classifiedErr.StatusCode, classifiedErr
	}

	providerResponse := GenerateResponse{}
	providerResponse.Text = response.Text()
	providerResponse.InputTokens = int64(response.UsageMetadata.PromptTokenCount)
	providerResponse.OutputTokens = int64(response.UsageMetadata.CandidatesTokenCount)
	providerResponse.TotalTokens = int64(response.UsageMetadata.TotalTokenCount)
	providerResponse.FinishReason = string(response.Candidates[0].FinishReason)
	return &providerResponse, 200, nil
}

func (g *GeminiProvider) GenerateStream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, <-chan error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		30*time.Second,
	)

	chunks := make(chan StreamChunk)
	errs := make(chan error, 1)

	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(input.Temperature)),
	}

	if input.MaxOutputTokens > 0 {
		config.MaxOutputTokens = int32(input.MaxOutputTokens)
	}

	go func() {
		defer cancel()
		defer close(chunks)
		defer close(errs)
		defer func() {
			if r := recover(); r != nil {
				errs <- ClassifyError("gemini", fmt.Errorf("gemini stream: recovered from panic: %v", r))
			}
		}()

		streamChunks := g.Client.Models.GenerateContentStream(
			ctx,
			input.Model,
			genai.Text(input.Prompt),
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
