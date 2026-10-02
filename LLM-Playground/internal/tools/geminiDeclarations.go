package tools

import "google.golang.org/genai"

func (tr *ToolRegistry) GeminiDeclarations() []*genai.FunctionDeclaration {
	declarations := make([]*genai.FunctionDeclaration, 0, len(tr.toolsByName))

	for _, tool := range tr.toolsByName {
		declaration := &genai.FunctionDeclaration{
			Name:                 tool.ToolName,
			Description:          tool.ToolDescription,
			ParametersJsonSchema: tool.ToolParameters,
		}

		declarations = append(declarations, declaration)
	}

	return declarations
}
