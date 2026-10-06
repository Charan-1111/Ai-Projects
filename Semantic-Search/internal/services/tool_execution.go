package services

import (
	"context"
	"fmt"
	"semantic-search/internal/models"
	"semantic-search/internal/tools"

	"github.com/bytedance/sonic"
)

func (s *Service) ExecuteTools(ctx context.Context, req tools.ExecuteToolRequest) (tools.ToolExecutionResponse, error) {
	switch req.ToolName {
	case "search_documents":
		var args models.SearchRequest
		if err := sonic.Unmarshal(req.Arguments, &args); err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		documentsResponse, err := s.SearchDocuments(ctx, args)
		if err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		resultBytes, err := sonic.Marshal(documentsResponse)
		if err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		var result map[string]any
		if err := sonic.Unmarshal(resultBytes, &result); err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		return tools.ToolExecutionResponse{Result: result}, nil
	case "get_document":
		var args models.FetchDocument
		if err := sonic.Unmarshal(req.Arguments, &args); err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		document, err := s.FetchDocument(ctx, args.DocId)
		if err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		resultBytes, err := sonic.Marshal(document)
		if err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		var result map[string]any
		if err := sonic.Unmarshal(resultBytes, &result); err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		return tools.ToolExecutionResponse{Result: result}, nil
	default:
		return tools.ToolExecutionResponse{}, fmt.Errorf("unknown tool: %s", req.ToolName)
	}
}
