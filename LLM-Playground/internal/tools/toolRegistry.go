package tools

import "llm-playground/internal/clients"

type RegisteredTool struct {
	ToolName        string
	ToolDescription string
	ToolParameters  map[string]any
	ServiceUrl      string
}

type ToolRegistry struct {
	clients     *clients.Clients
	apiFactory  clients.ApiFactory
	tools       map[string]string
	toolsByName map[string]RegisteredTool
}

func NewToolRegistry(clients *clients.Clients, apiFactory clients.ApiFactory, tools map[string]string) *ToolRegistry {
	return &ToolRegistry{
		clients:     clients,
		apiFactory:  apiFactory,
		tools:       tools,
		toolsByName: make(map[string]RegisteredTool),
	}
}
