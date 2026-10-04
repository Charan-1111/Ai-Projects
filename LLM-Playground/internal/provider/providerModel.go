package provider

import "llm-playground/internal/models"

type GenerateInput struct {
	SystemPrompt    string
	RequiredTools   []string
	Prompt          string
	History         []Message
	Model           string
	Temperature     float64
	MaxOutputTokens int64
}

type Message struct {
	Role    string
	Content string
}

type GenerateResponse struct {
	Text         string `json:"text"`
	InputTokens  int64  `json:"input_tokens"`
	OutputTokens int64  `json:"output_tokens"`
	TotalTokens  int64  `json:"total_tokens"`
	FinishReason string `json:"finish_reason"`
}

type StreamChunk struct {
	Delta        string
	InputTokens  int64
	OutputTokens int64
	FinishReason string
}

type GeneratorService struct {
	provider LLMProvider
}

func BuildGenerateInput(modelConfig models.ModelConfig, request *models.PromptRequest, globalInstructions string) (GenerateInput, error) {
	systemPrompt := globalInstructions

	if modelConfig.SystemInstructions != "" {
		systemPrompt += "\n\n" + modelConfig.SystemInstructions
	}

	input := GenerateInput{}
	input.MaxOutputTokens = modelConfig.MaxOutputTokens
	input.SystemPrompt = modelConfig.SystemInstructions
	input.RequiredTools = append([]string(nil), modelConfig.RequiredTools...)
	input.Model = modelConfig.ProviderModel
	input.Prompt = request.Prompt
	input.Temperature = request.Temperature

	return input, nil
}
