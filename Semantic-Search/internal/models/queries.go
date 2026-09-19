package models

type Queries struct {
	Create map[string]string `json:"create"`
	Fetch  Fetch             `json:"fetch"`
	Save   Save              `json:"save"`
	Edit   Edit              `json:"edit"`
	Delete Delete            `json:"delete"`
}

type Fetch struct {
	Messages         string `json:"messages"`
	SimilarDocuments string `json:"similar_documents"`
}

type Save struct {
	Embedding string `json:"embedding"`
}

type Edit struct {
	Document string `json:"document"`
}

type Delete struct {
	Document string `json:"document"`
}
