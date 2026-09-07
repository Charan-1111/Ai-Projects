package chat

import (
	"context"
	"fmt"
	"llm-playground/internal/provider"
	"llm-playground/internal/store/database"
)

type PersistentChatService struct {
	llm     provider.LLMProvider
	dbStore database.Repository
}

func NewPersistentChatService(llm provider.LLMProvider, dbStore database.Repository) *PersistentChatService {
	return &PersistentChatService{
		llm:     llm,
		dbStore: dbStore,
	}
}

func (pcs *PersistentChatService) Chat(ctx context.Context, requestId, conversationId string, request provider.GenerateInput) (*provider.GenerateResponse, error) {
	history, err := pcs.dbStore.Messages(ctx, conversationId)
	if err != nil {
		return nil, fmt.Errorf("load chat history: %w", err)
	}

	newConversation := len(history) == 0
	if newConversation {
		history = append(history, provider.Message{Role: "model", Content: "You are a helpful assistant."})
	}
	request.History = history

	response, _, err := pcs.llm.Generate(ctx, request)
	if err != nil {
		return nil, err
	}

	messages := []provider.Message{
		{Role: "user", Content: request.Prompt},
		{Role: "assistant", Content: response.Text},
	}
	if newConversation {
		messages = append([]provider.Message{history[0]}, messages...)
	}

	if err := pcs.dbStore.SaveMessages(ctx, requestId, conversationId, messages); err != nil {
		return nil, fmt.Errorf("save chat history: %w", err)
	}

	return response, nil
}
