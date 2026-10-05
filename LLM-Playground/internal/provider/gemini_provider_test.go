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
