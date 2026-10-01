package tools

type RegisteredTool struct {
	ToolName        string `json:"name"`
	ToolDescription string `json:"description"`
}
type ToolRegistry struct {
	ToolsByName map[string]RegisteredTool `json:"tools"`
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		ToolsByName: make(map[string]RegisteredTool),
	}
}
