package provider

import (
	"testing"
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
