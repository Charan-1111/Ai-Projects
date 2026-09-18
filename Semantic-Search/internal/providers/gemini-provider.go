package providers

import (
	"context"
	"fmt"
	"math"
	"os"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
}

func NewGeminiProvider(ctx context.Context) (*GeminiProvider, error) {
	// create gemini client here itself

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini provider api key is not present")
	}

	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Error creating gemini client : %w", err)
	}

	return &GeminiProvider{
		client: client,
	}, nil
}

func (gp *GeminiProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	outputDimensionality := int32(1536)
	result, err := gp.client.Models.EmbedContent(ctx, "gemini-embedding-001", genai.Text(text), &genai.EmbedContentConfig{
		OutputDimensionality: &outputDimensionality,
	})
	if err != nil {
		return []float32{}, err
	}

	return result.Embeddings[0].Values, nil
}

func (gp *GeminiProvider) CompareSimilarity(ctx context.Context, text1 string, text2 string) (float64, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(text1, genai.RoleUser),
		genai.NewContentFromText(text2, genai.RoleUser),
	}

	results, err := gp.client.Models.EmbedContent(
		ctx,
		"gemini-embedding-001",
		contents,
		&genai.EmbedContentConfig{
			TaskType: "SEMANTIC_SIMILARITY",
		},
	)

	if err != nil {
		return 0, err
	}

	vector1 := results.Embeddings[0].Values
	vector2 := results.Embeddings[1].Values

	similarity, err := cosineSimilarity(vector1, vector2)
	if err != nil {
		return 0, err
	}

	return similarity, nil
}

func cosineSimilarity(a, b []float32) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf(
			"vector lengths do not match: %d and %d",
			len(a),
			len(b),
		)
	}

	if len(a) == 0 {
		return 0, fmt.Errorf("vectors cannot be empty")
	}

	var dotProduct float64
	var magnitudeA float64
	var magnitudeB float64

	for i := range a {
		valueA := float64(a[i])
		valueB := float64(b[i])

		dotProduct += (valueA * valueB)
		magnitudeA += (valueA * valueA)
		magnitudeB += (valueB * valueB)
	}

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0, fmt.Errorf("cannot compare zero-magnitude vectors")
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB)), nil
}
