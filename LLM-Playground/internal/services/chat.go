package services

import (
	"context"
	"fmt"
	"llm-playground/internal/models"
	"llm-playground/internal/provider"
)

func (s *Services) Chat(ctx context.Context, requestId, conversationId string, request *models.PromptRequest) (*provider.GenerateResponse, error) {
	modelConfig, ok := s.resolveModelConfig(request)
	if !ok {
		return nil, fmt.Errorf("model configuration not found for request model %q and model_id %q", request.Model, request.ModelId)
	}

	input, _ := provider.BuildGenerateInput(modelConfig, request)

	// return s.inMemoryChatService.Chat(ctx, requestId, conversationId, input)
	return s.persistentChatService.Chat(ctx, requestId, conversationId, input)
}
