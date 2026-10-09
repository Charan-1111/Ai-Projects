package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/genai"
)

type MCPRegistry struct {
	sessions     []*mcp.ClientSession
	owners       map[string]*mcp.ClientSession
	declarations []*genai.FunctionDeclaration
}

func NewMCPRegistry() *MCPRegistry {
	return &MCPRegistry{
		owners: make(map[string]*mcp.ClientSession),
	}
}

// Call once during application startup.
func (r *MCPRegistry) Connect(
	ctx context.Context,
	serverName string,
	endpoint string,
) error {
	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "llm-playground",
			Version: "1.0.0",
		},
		nil,
	)

	session, err := client.Connect(
		ctx,
		&mcp.StreamableClientTransport{
			Endpoint: endpoint,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", serverName, err)
	}

	// Building a temporary list so failed discovery does not
	// leave partially registered tools.
	var discovered []*mcp.Tool
	params := &mcp.ListToolsParams{}

	for {
		result, err := session.ListTools(ctx, params)
		if err != nil {
			session.Close()
			return fmt.Errorf("discover tools from %s: %w", serverName, err)
		}

		discovered = append(discovered, result.Tools...)

		if result.NextCursor == "" {
			break
		}
		params.Cursor = result.NextCursor
	}

	names := make(map[string]bool)

	for _, tool := range discovered {
		if tool.Name == "" || names[tool.Name] {
			session.Close()
			return fmt.Errorf("invalid or duplicate tool name from %s", serverName)
		}
		if _, exists := r.owners[tool.Name]; exists {
			session.Close()
			return fmt.Errorf("tool name collision: %s", tool.Name)
		}
		names[tool.Name] = true
	}

	for _, tool := range discovered {
		r.owners[tool.Name] = session

		r.declarations = append(
			r.declarations,
			&genai.FunctionDeclaration{
				Name:                 tool.Name,
				Description:          tool.Description,
				ParametersJsonSchema: tool.InputSchema,
			},
		)
	}

	r.sessions = append(r.sessions, session)
	return nil
}

func (r *MCPRegistry) GeminiDeclarations() []*genai.FunctionDeclaration {
	return r.declarations
}
