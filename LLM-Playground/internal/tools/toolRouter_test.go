package tools

import (
	"context"
	"testing"

	"google.golang.org/genai"
)

type testToolExecutor struct {
	result map[string]any
}

func (e testToolExecutor) ToolExecution(context.Context, string, map[string]any) (map[string]any, error) {
	return e.result, nil
}

func TestToolRouterRoutesToRegisteredExecutor(t *testing.T) {
	router := NewToolRouter()
	restResult := map[string]any{"source": "rest"}
	mcpResult := map[string]any{"source": "mcp"}

	if err := router.Register([]*genai.FunctionDeclaration{{Name: "search_documents"}}, testToolExecutor{result: restResult}); err != nil {
		t.Fatal(err)
	}
	if err := router.Register([]*genai.FunctionDeclaration{{Name: "get_current_time"}}, testToolExecutor{result: mcpResult}); err != nil {
		t.Fatal(err)
	}

	result, err := router.ToolExecution(context.Background(), "search_documents", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["source"] != "rest" {
		t.Fatalf("expected REST executor result, got %v", result)
	}

	result, err = router.ToolExecution(context.Background(), "get_current_time", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["source"] != "mcp" {
		t.Fatalf("expected MCP executor result, got %v", result)
	}
}

func TestToolRouterUsesLastRegisteredToolName(t *testing.T) {
	router := NewToolRouter()
	restResult := map[string]any{"source": "rest"}
	mcpResult := map[string]any{"source": "mcp"}

	if err := router.Register([]*genai.FunctionDeclaration{{Name: "get_current_time"}}, testToolExecutor{result: restResult}); err != nil {
		t.Fatal(err)
	}
	if err := router.Register([]*genai.FunctionDeclaration{{Name: "get_current_time"}}, testToolExecutor{result: mcpResult}); err != nil {
		t.Fatal(err)
	}

	if declarations := router.GeminiDeclarations(); len(declarations) != 1 {
		t.Fatalf("expected one deduplicated declaration, got %d", len(declarations))
	}
	result, err := router.ToolExecution(context.Background(), "get_current_time", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["source"] != "mcp" {
		t.Fatalf("expected last registered executor result, got %v", result)
	}
}
