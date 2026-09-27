package models

type AskDocument struct {
	Prompt       string       `json:"prompt"`
	NoOfDocs     int          `json:"noOfDocs"`
	MinimumScore float64      `json:"minimumScore"`
	Filters      AskFilters   `json:"filters"`
}

type AskFilters struct {
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
}