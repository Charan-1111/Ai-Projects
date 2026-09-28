package services

import (
	"ask-my-documents/internal/clients"
	"ask-my-documents/internal/models"
	"ask-my-documents/internal/utils"
	"context"
	"fmt"
	"sort"
	"strings"
)

func (s *Service) AskDocuments(ctx context.Context, req models.AskDocument) (*models.AnswerResponse, error) {
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

	semanticClient := clients.NewSemanticSearchClient(s.apiFactory, s.config.ExternalApis.SemanticSearch, s.clients)
	semanticResponse, err := semanticClient.Search(ctx, req.Prompt, noOfDocs, minimumScore, filters)
	if err != nil {
		return nil, fmt.Errorf("external semantic search call failed: %w", err)
	}

	contextText, sources := utils.BuildContext(semanticResponse.Contents.Documents)
	if len(sources) == 0 {
		return &models.AnswerResponse{
			Answer: "I could not find that information in the uploaded documents.",
		}, nil
	}

	llmPrompt := utils.BuildPrompt(req.Prompt, contextText)
	llmClient := clients.NewLLMClient(s.apiFactory, s.config.ExternalApis.Chat, s.clients)
	llmResponse, err := llmClient.Generate(ctx, llmPrompt)
	if err != nil {
		return nil, fmt.Errorf("llm generation failed: %w", err)
	}
	if strings.TrimSpace(llmResponse.Text) == "" {
		return nil, fmt.Errorf("llm response is empty")
	}

	citations := make([]models.Citation, 0, len(sources))
	keys := make([]int, 0, len(sources))
	for ref := range sources {
		keys = append(keys, ref)
	}
	sort.Ints(keys)
	for _, ref := range keys {
		if !strings.Contains(llmResponse.Text, fmt.Sprintf("[%d]", ref)) {
			continue
		}

		source := sources[ref]
		citations = append(citations, models.Citation{
			Reference:  source.Reference,
			DocumentID: source.DocumentId,
			ChunkID:    source.ChunkId,
			Title:      source.Title,
		})
	}

	return &models.AnswerResponse{
		Answer:    llmResponse.Text,
		Citations: citations,
	}, nil
}
