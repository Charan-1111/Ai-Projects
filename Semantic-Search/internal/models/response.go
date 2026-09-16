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
