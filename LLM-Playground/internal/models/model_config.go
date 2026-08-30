package models

type ModelConfig struct {
	DisplayName                string  `json:"display_name"`
	ProviderModel              string  `json:"provider_model"`
	InputCostPerMillionTokens  float64 `json:"input_cost_per_million_tokens"`
	OutputCostPerMillionTokens float64 `json:"output_cost_per_million_tokens"`
	MaxOutputTokens            int64   `json:"max_output_tokens"`
	SupportsStreaming          bool    `json:"supports_streaming"`
}
