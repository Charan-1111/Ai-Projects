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

func (g *GeminiProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResponse, error) {
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
		return &GenerateResponse{}, fmt.Errorf("Failed to generate response : %w", err)
	}

	if response == nil || response.UsageMetadata == nil {
		return &GenerateResponse{}, fmt.Errorf("Response generation error")
	}

	providerResponse := GenerateResponse{}
	providerResponse.Text = response.Text()
	providerResponse.InputTokens = int64(response.UsageMetadata.PromptTokenCount)
	providerResponse.OutputTokens = int64(response.UsageMetadata.CandidatesTokenCount)
	providerResponse.TotalTokens = int64(response.UsageMetadata.TotalTokenCount)
	providerResponse.FinishReason = "stop"
	return &providerResponse, nil
}
