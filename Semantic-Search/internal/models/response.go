package models

type DocumentResponse struct {
	DocId           string `json:"docId"`
	DocTitle        string `json:"docTitle"`
	EmbeddingStatus string `json:"embeddingStatus"`
}

type EmbedResponse struct {
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
}

type SimilarDocuments struct {
	Id          string  `json:"id"`
	DocTitle    string  `json:"title"`
	DocContent  string  `json:"content"`
	DocCategory string  `json:"category"`
	Similarity  float32 `json:"similarity"`
}

type SearchResponse struct {
	Query            string             `json:"query"`
	SimilarDocuments []SimilarDocuments `json:"documentResults"`
	ResultCount      int                `json:"resultCount"`
	DurationMs       int                `json:"durationMs"`
}

type ChunkDetails struct {
	Id            string         `json:"id"`
	DocumentId    string         `json:"document_id"`
	Content       string         `json:"content"`
	Similarity    float32        `json:"similarity"`
	ChunkIndex    int            `json:"chunk_index"`
	DocumentTitle string         `json:"title"`
	Category      string         `json:"category"`
	Metadata      map[string]any `json:"metadata"`
}

type ChunkedResponse struct {
	Query       string         `json:"query"`
	Documents   []ChunkDetails `json:"documents"`
	ResultCount int            `json:"resultCount"`
	DurationMs  int            `json:"durationMs"`
}
