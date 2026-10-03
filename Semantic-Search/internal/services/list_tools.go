package services

import (
	"context"
	"semantic-search/internal/tools"
)

func (s *Service) ListTools(ctx context.Context) tools.ListToolsResponse {
	registeredTools := make([]tools.ToolDefinition, 0)

	registeredTools = append(registeredTools, tools.SearchDocumentsDefinition())

	return tools.ListToolsResponse{
		Tools: registeredTools,
	}
}
