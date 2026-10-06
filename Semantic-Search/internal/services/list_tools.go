package services

import (
	"context"
	"semantic-search/internal/tools"
)

func (s *Service) ListTools(ctx context.Context) tools.ListToolsResponse {
	registeredTools := make([]tools.ToolDefinition, 0)

	registeredTools = append(registeredTools, tools.SearchDocumentsDefinition())
	registeredTools= append(registeredTools, tools.FetchDocumentDefinition())
	
	return tools.ListToolsResponse{
		Tools: registeredTools,
	}
}
