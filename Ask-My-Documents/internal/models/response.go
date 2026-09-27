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
