package models

type PromptRequest struct {
	Prompt          string  `json:"prompt"`
	Model           string  `json:"model"`
	ModelId         string  `json:"model_id"`
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int64   `json:"max_output_tokens"`
	Stream          bool    `json:"stream"`
}
