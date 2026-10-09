package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *MCPRegistry) ToolExecution(ctx context.Context, toolName string, args map[string]any) (map[string]any, error) {
	session, exists := r.owners[toolName]
	if !exists {
		return nil, fmt.Errorf("unknown tool %q", toolName)
	}

	if args == nil {
		args = map[string]any{}
	}

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	})
	if err != nil {
		return nil, fmt.Errorf("call tool %q: %w", toolName, err)
	}

	// passing text results through, including failures
	texts := make([]string, 0)
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			texts = append(texts, text.Text)
		}
	}

	if result.IsError {
		return map[string]any{
			"error":    true,
			"messages": texts,
		}, nil
	}

	return map[string]any{
		"result":   result.StructuredContent,
		"messages": texts,
	}, nil
}

func (r *MCPRegistry) Close() {
	for _, session := range r.sessions {
		session.Close()
	}
}
