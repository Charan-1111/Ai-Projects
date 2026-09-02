package provider

import "context"

type LLMProvider interface {
	Generate(ctx context.Context, input GenerateInput) (*GenerateResponse, error)
	GenerateStream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, <-chan error)
}
