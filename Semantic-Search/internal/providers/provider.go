package providers

import "context"

type LLMProvider interface {
	EmbedText(ctx context.Context, text string) ([]float32, error)
	CompareSimilarity(ctx context.Context, text1 string, text2 string) (float64, error)
}
