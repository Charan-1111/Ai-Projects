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

type ModelResponse struct {
	RequestId         string  `json:"request_id"`
	Model             string  `json:"model"`
	Response          string  `json:"response"`
	Usage             Usage   `json:"usage"`
	LatencyMs         int64   `json:"latency_ms"`
	EstimatedCostUsed float64 `json:"estimated_cost_used"`
	FinishReason      string  `json:"finish_reason"`
}

type Usage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

type ToolResponse struct {
	Tools []ToolDetails `json:"tools"`
}

type ToolDetails struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ToolExecutionResponse struct {
	Result map[string]any      `json:"result,omitempty"`
	Error  *ToolExecutionError `json:"error,omitempty"`
}

type ToolExecutionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
