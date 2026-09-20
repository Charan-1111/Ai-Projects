package models

type Queries struct {
	Create map[string]string `json:"create"`
	Fetch  Fetch             `json:"fetch"`
	Save   Save              `json:"save"`
	Edit   Edit              `json:"edit"`
	Delete Delete            `json:"delete"`
}

type Fetch struct {
	Messages               string `json:"messages"`
	SimilarDocuments       string `json:"similar_documents"`
	ChunkedDocuments       string `json:"chunked_documents"`
	FilterChunkedDocuments string `json:"filter_chunked_documents"`
}

type Save struct {
	Embedding string `json:"embedding"`
	Chunks    string `json:"chunks"`
}

type Edit struct {
	Document       string `json:"document"`
	IndexingStatus string `json:"index_status"`
}

type Delete struct {
	Document string `json:"document"`
}
