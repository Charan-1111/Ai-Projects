package services

import (
	"ask-my-documents/internal/models"
	"ask-my-documents/internal/utils"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
)

func (s *Service) AskDocuments(ctx context.Context, req models.AskDocument) (models.SemanticResponse, error) {
	var semanticResponse models.SemanticResponse

	noOfDocs := req.NoOfDocs
	if noOfDocs <= 0 {
		noOfDocs = 5
	}

	minimumScore := req.MinimumScore
	if minimumScore <= 0 {
		minimumScore = 0.7
	}

	filters := map[string]string{}
	if req.Filters.Category != "" {
		filters["category"] = req.Filters.Category
	}
	if req.Filters.Difficulty != "" {
		filters["difficulty"] = req.Filters.Difficulty
	}

	apiClient := s.apiFactory.Create(
		s.config.ExternalApis.SemanticSearch,
		"POST",
		map[string]string{},
		map[string]any{
			"query":        req.Prompt,
			"noOfDocs":     noOfDocs,
			"minimumScore": minimumScore,
			"filters":      filters,
		},
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		s.clients,
	)

	respDataBytes, err := apiClient.ApiCall(ctx)
	if err != nil {
		return semanticResponse, fmt.Errorf("external semantic search call failed: %w", err)
	}

	if err = sonic.Unmarshal(respDataBytes, &semanticResponse); err != nil {
		return semanticResponse, fmt.Errorf("unmarshal semantic search response: %w", err)
	}

	contextText, sources := utils.BuildContext(semanticResponse.Contents.Documents)

	if len(sources) == 0 {
		semanticResponse = models.SemanticResponse{
			Code:    0,
			Message: "I could not find that information in the uploaded documents.",
		}
		return semanticResponse, nil
	}

	llmPrompt := utils.BuildPrompt(req.Prompt, contextText)

	fmt.Println(llmPrompt)
	return semanticResponse, nil
}
