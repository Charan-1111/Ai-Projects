package tools

import "encoding/json"

type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ListToolsResponse struct {
	Tools []ToolDefinition `json:"tools"`
}

type ExecuteToolRequest struct {
	ToolName      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type ToolError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolExecutionResponse struct {
	Result map[string]any `json:"result,omitempty"`
	Error  *ToolError      `json:"error,omitempty"`
}

