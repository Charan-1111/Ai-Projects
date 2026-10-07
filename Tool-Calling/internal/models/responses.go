package models

type GetTimeResponse struct {
	Answer        string   `json:"answer"`
	ToolsExecutes []string `json:"tools_executed"`
}

type TimeOutput struct {
	Timezone string `json:"timezone"`
	Datetime string `json:"datetime"`
}
