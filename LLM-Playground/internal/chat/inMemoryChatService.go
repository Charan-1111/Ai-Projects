package chat

import (
	"context"
	"llm-playground/internal/provider"
	"sync"
)

type InMemoryChatService struct {
	llm      provider.LLMProvider
	mu       sync.Mutex
	messages map[string][]provider.Message
}

func NewInMemoryChatService(llm provider.LLMProvider) *InMemoryChatService {
	return &InMemoryChatService{
		llm:      llm,
		messages: make(map[string][]provider.Message),
	}
}

func (ics *InMemoryChatService) Chat(ctx context.Context, requestId, conversationId string, request provider.GenerateInput) (*provider.GenerateResponse, error) {
	ics.mu.Lock()
	defer ics.mu.Unlock()

	history := append([]provider.Message(nil), ics.messages[conversationId]...)
	request.History = history
	response, _, err := ics.llm.Generate(ctx, request)
	if err != nil {
		return nil, err
	}

	history = append(history,
		provider.Message{Role: "user", Content: request.Prompt},
		provider.Message{Role: "assistant", Content: response.Text},
	)
	ics.messages[conversationId] = history

	return response, nil
}
