package providers

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
}

func NewGeminiProvider(ctx context.Context) (*GeminiProvider, error) {
	// create gemini client here itself

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini provider api key is not present")
	}

	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Error creating gemini client : %w", err)
	}

	return &GeminiProvider{
		client: client,
	}, nil
}
