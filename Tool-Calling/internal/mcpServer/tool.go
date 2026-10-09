package mcpServer

import (
	"context"
	"tool-calling/internal/models"
	"tool-calling/internal/services"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTools(server *mcp.Server, svc *services.Service) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_current_time",
			Description: "Get the current time in an IANA timezone.",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input models.TimeInput,
		) (*mcp.CallToolResult, models.TimeOutput, error) {
			output, err := svc.GetCurrentTime(ctx, input)
			return nil, output, err
		},
	)
}
