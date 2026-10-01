package services

import (
	"context"
	"fmt"
	"tool-calling/internal/tools"

	"github.com/bytedance/sonic"
)

func (s *Service) ToolExecution(ctx context.Context, req tools.ExecuteToolRequest) (tools.ToolExecutionResponse, error) {
	switch req.ToolName {
	case "get_current_time":
		// Execute the get_current_time tool logic
		var args tools.CurrentTimeArgs
		if err := sonic.Unmarshal(req.Arguments, &args); err != nil {
			// Handle error
			return tools.ToolExecutionResponse{}, err
		}

		result, err := tools.GetCurrentTime(args)
		if err != nil {
			return tools.ToolExecutionResponse{}, err
		}

		return tools.ToolExecutionResponse{
			Result: map[string]any{
				"timezone": result.Timezone,
				"time":     result.CurrentTime,
			},
		}, nil
	default:
		// Handle unknown tool name
		return tools.ToolExecutionResponse{}, fmt.Errorf("unknown tool: %s", req.ToolName)
	}
}
