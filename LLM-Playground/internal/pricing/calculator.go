package pricing

import (
	"llm-playground/internal/models"
	"llm-playground/internal/provider"
)

func PriceCalculator(model models.ModelConfig, response *provider.GenerateResponse) (models.Usage, float64) {
	usage := models.Usage{}
	
	usage.InputTokens = int64(response.InputTokens)
	usage.OutputTokens = int64(response.OutputTokens)
	usage.TotalTokens = int64(response.TotalTokens)

	const tokensPerMillion = 1_000_000.0

	totalCost := (float64(usage.InputTokens) / tokensPerMillion * model.InputCostPerMillionTokens) +
		(float64(usage.OutputTokens) / tokensPerMillion * model.OutputCostPerMillionTokens)

	return usage, totalCost
}
