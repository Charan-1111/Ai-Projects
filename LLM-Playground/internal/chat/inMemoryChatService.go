package chat

import (
	"llm-playground/internal/provider"
	"sync"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type InMemoryChatService struct {
	llm provider.LLMProvider
	my  sync.Mutex
	messages map[string][]Message
}

func NewInMemoryChatService(llm provider.LLMProvider) *InMemoryChatService {
	return &InMemoryChatService{
		llm:      llm,
		messages: make(map[string][]Message),
	}
}
