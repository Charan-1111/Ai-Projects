package mcpServer

import (
	"net/http"
	"tool-calling/internal/services"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewHandler(svc *services.Service) http.Handler {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "tool-calling",
			Version: "1.0.0",
		},
		nil,
	)

	// Tool registering
	registerTools(server, svc)

	return mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{
			Stateless:    true,
			JSONResponse: true,
		},
	)
}
