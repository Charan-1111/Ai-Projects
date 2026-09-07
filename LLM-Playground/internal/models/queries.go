package models

type Queries struct {
	Create map[string]string `json:"create"`
	Fetch  Fetch             `json:"fetch"`
	Save   Save              `json:"save"`
}

type Fetch struct {
	Messages string `json:"messages"`
}

type Save struct {
	Conversation string `json:"conversation"`
	Message      string `json:"message"`
}
