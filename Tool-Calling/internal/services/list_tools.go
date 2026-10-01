package services

import "tool-calling/internal/tools"

func (s *Service) ListTools() tools.ListToolsResponse {
	toolList := make([]tools.ToolDefinition, 0)

	toolList = append(toolList, tools.CurrentTimeDefiniation())

	return tools.ListToolsResponse{
		Tools: toolList,
	}
}
