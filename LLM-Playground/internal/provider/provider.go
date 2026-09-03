package provider

import "context"

type StreamError struct {
	StatusCode int
	Err        error
}

func (e *StreamError) Error() string {
	return e.Err.Error()
}

func (e *StreamError) Unwrap() error {
	return e.Err
}

type LLMProvider interface {
	Generate(ctx context.Context, input GenerateInput) (*GenerateResponse, int, error)
	GenerateStream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, <-chan error)
}
