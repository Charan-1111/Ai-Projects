package models

type SemanticResponse struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	Contents SemanticContents `json:"contents"`
}

type SemanticContents struct {
	Query       string              `json:"query"`
	Documents   []SemanticDocuments `json:"documents"`
	ResultCount int                 `json:"resultCount"`
	DurationMs  int                 `json:"durationMs"`
}

type SemanticDocuments struct {
	Id         string         `json:"id"`
	DocumentId string         `json:"document_id"`
	Content    string         `json:"content"`
	Similarity float32        `json:"similarity"`
	Title      string         `json:"title"`
	Category   string         `json:"category"`
	Metadata   map[string]any `json:"metadata"`
}

type LLMResponse struct {
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

type AnswerResponse struct {
	Answer    string     `json:"answer"`
	Citations []Citation `json:"citations"`
}

type Citation struct {
	Reference  int    `json:"reference"`
	DocumentID string `json:"document_id"`
	ChunkID    string `json:"chunk_id"`
	Title      string `json:"title"`
}
