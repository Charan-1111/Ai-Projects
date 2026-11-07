package provider

import (
	"testing"

	"google.golang.org/genai"
)

func TestGeminiContentsSkipsEmptyHistoryMessages(t *testing.T) {
	contents := geminiContents(GenerateInput{
		Prompt: "What time is it?",
		History: []Message{
			{Role: "user", Content: "Earlier question"},
			{Role: "model", Content: ""},
		},
	})

	if len(contents) != 2 {
		t.Fatalf("expected two content entries, got %d", len(contents))
	}
	for index, content := range contents {
		if len(content.Parts) == 0 {
			t.Fatalf("content %d has no parts", index)
		}
		for partIndex, part := range content.Parts {
			if part == nil || part.Text == "" {
				t.Fatalf("content %d part %d is empty", index, partIndex)
			}
		}
	}
}

func TestGeminiFunctionCallsCollectsEveryCall(t *testing.T) {
	content := &genai.Content{
		Parts: []*genai.Part{
			{FunctionCall: &genai.FunctionCall{Name: "get_current_time"}},
			{FunctionCall: &genai.FunctionCall{Name: "search_goroutines"}},
		},
	}

	calls := geminiFunctionCalls(content)
	if len(calls) != 2 {
		t.Fatalf("expected two function calls, got %d", len(calls))
	}
	if calls[0].Name != "get_current_time" || calls[1].Name != "search_goroutines" {
		t.Fatalf("unexpected function calls: %q, %q", calls[0].Name, calls[1].Name)
	}
}

func TestGeminiRequiredToolsOnlyConstrainInitialRequest(t *testing.T) {
	config := geminiGenerateConfig(
		GenerateInput{RequiredTools: []string{"search_documents"}},
		[]*genai.FunctionDeclaration{{Name: "search_documents"}, {Name: "get_document"}},
	)
	functionCallingConfig := config.ToolConfig.FunctionCallingConfig
	if functionCallingConfig.Mode != genai.FunctionCallingConfigModeAny {
		t.Fatalf("expected initial mode ANY, got %q", functionCallingConfig.Mode)
	}
	if len(functionCallingConfig.AllowedFunctionNames) != 1 || functionCallingConfig.AllowedFunctionNames[0] != "search_documents" {
		t.Fatalf("expected only search_documents to be required initially, got %v", functionCallingConfig.AllowedFunctionNames)
	}

	config.ToolConfig = geminiAutomaticToolConfig()
	functionCallingConfig = config.ToolConfig.FunctionCallingConfig
	if functionCallingConfig.Mode != genai.FunctionCallingConfigModeAuto {
		t.Fatalf("expected subsequent mode AUTO, got %q", functionCallingConfig.Mode)
	}
	if len(functionCallingConfig.AllowedFunctionNames) != 0 {
		t.Fatalf("expected subsequent requests to allow every declared tool, got %v", functionCallingConfig.AllowedFunctionNames)
	}
}
