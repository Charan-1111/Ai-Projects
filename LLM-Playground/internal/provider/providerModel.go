package provider

import "llm-playground/internal/models"

type GenerateInput struct {
	SystemPrompt    string
	Prompt          string
	Model           string
	Temperature     float64
	MaxOutputTokens int64
}

type GenerateResponse struct {
	Text         string
	InputTokens  int64
	OutputTokens int64
	TotalTokens  int64
	FinishReason string
}

type StreamChunk struct {
	Delta           string
	InputTokens     int64
	OutputTokens    int64
	FinishedReasong string
}

type GeneratorService struct {
	provider LLMProvider
}

func BuildGenerateInput(modelConfig models.ModelConfig, request *models.PromptRequest) (GenerateInput, error) {
	input := GenerateInput{}
	input.MaxOutputTokens = modelConfig.MaxOutputTokens
	input.Model = modelConfig.ProviderModel
	input.Prompt = request.Prompt
	input.Temperature = request.Temperature

	return input, nil
}
