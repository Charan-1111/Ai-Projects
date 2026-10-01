package tools

func CurrentTimeDefiniation() ToolDefinition{
	return ToolDefinition{
		Name: "get_current_time",
		Description: "Get the current date and time in an IANA timezone. " +
			"Use for questions about the current time. " +
			"Ask for clarification when the location is ambiguous.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"timezone": map[string]any{
					"type": "string",
					"description": "IANA timezone, such as " +
						"Asia/Kolkata or America/New_York",
				},
			},
			"required": []string{"timezone"},
		},
	}
}