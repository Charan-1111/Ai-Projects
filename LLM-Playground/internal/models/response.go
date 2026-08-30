package models

type AvailableModels struct {
	Id               string `json:"id"`
	DisplayName      string `json:"display_name"`
	MaxOutputTokens  int64  `json:"max_output_tokens"`
	SupportStreaming bool   `json:"supports_streaming"`
}

type AvailableModelsResponse struct {
	DefaultModel    string            `json:"default_model"`
	AvailableModels []AvailableModels `json:"models_available"`
}
