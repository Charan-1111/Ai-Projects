package clients

import (
	"ask-my-documents/internal/models"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
)

type LLMClient struct {
	apiFactory ApiFactory
	url        string
	httpClient *Clients
}

func NewLLMClient(apiFactory ApiFactory, url string, httpClient *Clients) *LLMClient {
	return &LLMClient{
		apiFactory: apiFactory,
		url:        url,
		httpClient: httpClient,
	}
}

func (c *LLMClient) Generate(ctx context.Context, prompt string) (*models.LLMResponse, error) {
	apiClient := c.apiFactory.Create(
		c.url,
		"POST",
		map[string]string{},
		map[string]any{
			"prompt":            prompt,
			"model":             "gemini-3.5-flash-lite",
			"model_id":          "fast-model",
			"temperature":       0.6,
			"max_output_tokens": 10000,
			"stream":            false,
		},
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		c.httpClient,
	)

	resp, err := apiClient.ApiCall(ctx)
	if err != nil {
		return nil, fmt.Errorf("llm request failed: %w", err)
	}

	var llmResponse models.LLMResponse
	err = sonic.Unmarshal(resp, &llmResponse)
	if err != nil {
		return nil, err
	}

	return &llmResponse, nil
}
