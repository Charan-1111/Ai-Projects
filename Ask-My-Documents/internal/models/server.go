package models

type Server struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type ExternalApis struct {
	SemanticSearch string `json:"semanticSearch"`
}
