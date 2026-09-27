package clients

import (
	"ask-my-documents/internal/models"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
)

type SemanticSearchClient struct {
	apiFactory ApiFactory
	url        string
	httpClient *Clients
}

func NewSemanticSearchClient(apiFactory ApiFactory, url string, httpClient *Clients) *SemanticSearchClient {
	return &SemanticSearchClient{
		apiFactory: apiFactory,
		url:        url,
		httpClient: httpClient,
	}
}

func (c *SemanticSearchClient) Search(ctx context.Context, prompt string, noOfDocs int, minimumScore float64, filters map[string]string) (models.SemanticResponse, error) {
	var result models.SemanticResponse

	apiClient := c.apiFactory.Create(
		c.url,
		"POST",
		map[string]string{},
		map[string]any{
			"query":        prompt,
			"noOfDocs":     noOfDocs,
			"minimumScore": minimumScore,
			"filters":      filters,
		},
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		c.httpClient,
	)

	respDataBytes, err := apiClient.ApiCall(ctx)
	if err != nil {
		return result, fmt.Errorf("semantic search http request failed: %w", err)
	}

	if err = sonic.Unmarshal(respDataBytes, &result); err != nil {
		return result, fmt.Errorf("unmarshal semantic search response: %w", err)
	}

	return result, nil
}
