package tools

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type toolExecutor interface {
	ToolExecution(ctx context.Context, toolName string, args map[string]any) (map[string]any, error)
}

type ToolRouter struct {
	declarations []*genai.FunctionDeclaration
	executors    map[string]toolExecutor
	positions    map[string]int
}

func NewToolRouter() *ToolRouter {
	return &ToolRouter{
		executors: make(map[string]toolExecutor),
		positions: make(map[string]int),
	}
}

func (r *ToolRouter) Register(declarations []*genai.FunctionDeclaration, executor toolExecutor) error {
	for _, declaration := range declarations {
		if declaration == nil || declaration.Name == "" {
			return fmt.Errorf("cannot register tool with an empty name")
		}
	}

	for _, declaration := range declarations {
		if position, exists := r.positions[declaration.Name]; exists {
			r.declarations[position] = declaration
		} else {
			r.positions[declaration.Name] = len(r.declarations)
			r.declarations = append(r.declarations, declaration)
		}
		r.executors[declaration.Name] = executor
	}
	return nil
}

func (r *ToolRouter) GeminiDeclarations() []*genai.FunctionDeclaration {
	return r.declarations
}

func (r *ToolRouter) ToolExecution(ctx context.Context, toolName string, args map[string]any) (map[string]any, error) {
	executor, exists := r.executors[toolName]
	if !exists {
		return nil, fmt.Errorf("unknown tool %q", toolName)
	}
	return executor.ToolExecution(ctx, toolName, args)
}
