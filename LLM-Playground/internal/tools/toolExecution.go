package tools

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"net/url"

	"github.com/bytedance/sonic"
)

func (tr *ToolRegistry) ToolExecution(ctx context.Context, toolName string, args map[string]any) (map[string]any, error) {
	tool, exists := tr.toolsByName[toolName]
	if !exists {
		return nil, fmt.Errorf("unknown tool %q", toolName)
	}
	if args == nil {
		args = map[string]any{}
	}

	arguments, err := sonic.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal arguments for tool %q: %w", toolName, err)
	}

	body, err := sonic.Marshal(models.ExecuteToolRequest{
		ToolName:  toolName,
		Arguments: arguments,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request for tool %q: %w", toolName, err)
	}
	endpoint, err := url.JoinPath(tool.ServiceUrl, "tools", "execute")
	if err != nil {
		return nil, fmt.Errorf("build execution URL for tool %q: %w", toolName, err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	toolApi := tr.apiFactory.Create(
		endpoint,
		"POST",
		map[string]string{},
		body,
		headers,
		tr.clients,
	)

	toolResponse, callErr := toolApi.ApiCall(ctx)
	if len(toolResponse) == 0 {
		if callErr != nil {
			return nil, fmt.Errorf("call tool %q: %w", toolName, callErr)
		}
		return nil, fmt.Errorf("tool %q returned an empty response", toolName)
	}

	var decoded models.ToolExecutionResponse
	if err := sonic.Unmarshal(toolResponse, &decoded); err != nil {
		if callErr != nil {
			return nil, fmt.Errorf("call tool %q: %w (decode response: %v)", toolName, callErr, err)
		}
		return nil, fmt.Errorf("decode response from tool %q: %w", toolName, err)
	}
	if decoded.Error != nil {
		return nil, fmt.Errorf("tool %q failed (%d): %s", toolName, decoded.Error.Code, decoded.Error.Message)
	}
	if callErr != nil {
		return nil, fmt.Errorf("call tool %q: %w", toolName, callErr)
	}
	return decoded.Result, nil
}
