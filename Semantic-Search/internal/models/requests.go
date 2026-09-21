package models

import "time"

type Compare struct {
	Original string `json:"original"`
	Compare  string `json:"compare"`
}

type EmbedRequest struct {
	Text string `json:"text"`
}

type Document struct {
	DocId       string         `json:"docId"`
	DocTitle    string         `json:"docTitle"`
	DocContent  string         `json:"docContent"`
	DocCategory string         `json:"docCategory"`
	DocSource   string         `json:"docSource"`
	MetaData    map[string]any `json:"metaData"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type SearchRequest struct {
	Query        string  `json:"query"`
	NoofDocs     int     `json:"noOfDocs"`
	Filters      Filters `json:"filters"`
	MinimumScore float32 `json:"minimumScore"`
}

type Filters struct {
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
}
