package tools

import (
	"context"
	"fmt"
	"llm-playground/internal/models"

	"github.com/bytedance/sonic"
)

func (tr *ToolRegistry) RegisterTools() {
	// wg := &sync.WaitGroup{}

	for _, toolUrl := range tr.tools {
		finalUrl := toolUrl + "tools"

		apiCall := tr.apiFactory.Create(finalUrl, "GET", nil, nil, nil, tr.clients)

		apiResponse, err := apiCall.ApiCall(context.Background())
		if err != nil {
			fmt.Printf("Error occurred while calling API for tool %s: %v\n", toolUrl, err)
			continue
		}

		var toolData models.ToolResponse

		err = sonic.Unmarshal(apiResponse, &toolData)
		if err != nil {
			fmt.Printf("Error occurred while unmarshaling API response for tool %s: %v\n", toolUrl, err)
			continue
		}

		for _, tool := range toolData.Tools {
			if tool.Name == "" {
				fmt.Printf("Tool name is empty for tool from URL %s. Skipping registration.\n", toolUrl)
				continue
			}

			if _, exists := tr.toolsByName[tool.Name]; exists {
				fmt.Printf("Tool with name %s already exists in the registry. Skipping registration.\n", tool.Name)
				continue
			}

			tr.toolsByName[tool.Name] = RegisteredTool{
				ToolName:        tool.Name,
				ToolDescription: tool.Description,
				ToolParameters:  tool.Parameters,
				ServiceUrl:      toolUrl,
			}
		}
	}
}
