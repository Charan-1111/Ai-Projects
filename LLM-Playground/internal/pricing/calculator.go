package pricing

import (
	"llm-playground/internal/models"

	"google.golang.org/genai"
)

func PriceCalculator(model models.ModelConfig, response *genai.GenerateContentResponse) (models.Usage, float64) {
	usage := models.Usage{}

	if response == nil || response.UsageMetadata == nil {
		return usage, 0
	}

	usage.InputTokens = int64(response.UsageMetadata.PromptTokenCount)
	usage.OutputTokens = int64(response.UsageMetadata.CandidatesTokenCount)
	usage.TotalTokens = int64(response.UsageMetadata.TotalTokenCount)

	const tokensPerMillion = 1_000_000.0

	totalCost := (float64(usage.InputTokens) / tokensPerMillion * model.InputCostPerMillionTokens) +
		(float64(usage.OutputTokens) / tokensPerMillion * model.OutputCostPerMillionTokens)

	return usage, totalCost
}
