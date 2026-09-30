package models

type GetTimeResponse struct {
	Answer        string   `json:"answer"`
	ToolsExecutes []string `json:"tools_executed"`
}
