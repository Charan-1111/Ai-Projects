package tools

func SearchDocumentsDefinition() ToolDefinition {
	return ToolDefinition{
		Name:        "search_documents",
		Description: "Search the vector database for semantically relevant documents based on a natural-language query.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Natural-language question or search phrase to find relevant documents.",
				},
				"noOfDocs": map[string]any{
					"type":        "integer",
					"description": "Number of matching documents to return. Defaults to 10 when zero or negative.",
				},
				"minimumScore": map[string]any{
					"type":        "number",
					"description": "Optional minimum similarity score to apply when filtering chunk results.",
				},
				"filters": map[string]any{
					"type": "object",
					"description": "Optional filters to narrow document results.",
					"properties": map[string]any{
						"category": map[string]any{
							"type":        "string",
							"description": "Optional category to filter by.",
						},
						"difficulty": map[string]any{
							"type":        "string",
							"description": "Optional difficulty to filter by.",
						},
					},
					"additionalProperties": false,
				},
			},
			"required": []string{"query"},
		},
	}
}