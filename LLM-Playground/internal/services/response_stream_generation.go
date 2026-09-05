package services

import (
	"context"
	"errors"
	"fmt"
	"llm-playground/internal/models"
	"llm-playground/internal/provider"
	"llm-playground/internal/utils"
	"math/rand/v2"
	"time"
)

func (s *Services) ResponseStreamGeneration(ctx context.Context, request *models.PromptRequest) (<-chan provider.StreamChunk, <-chan error, error) {
	if request == nil {
		return nil, nil, fmt.Errorf("request body is required")
	}

	modelConfig, ok := s.resolveModelConfig(request)
	if !ok {
		return nil, nil, fmt.Errorf("model configuration not found for request model %q and model_id %q", request.Model, request.ModelId)
	}
	if s.config.Retries.MaxAttempts <= 0 || s.config.Retries.InitialDelayMs < 0 || s.config.Retries.MaximumDelayMs < 0 || s.config.Retries.MaximumDelayMs < s.config.Retries.InitialDelayMs {
		return nil, nil, fmt.Errorf("invalid retry configuration")
	}

	input, _ := provider.BuildGenerateInput(modelConfig, request)

	outChunks := make(chan provider.StreamChunk)
	outErrs := make(chan error, 1)

	go func() {
		defer close(outChunks)
		defer close(outErrs)

		delay := time.Duration(s.config.Retries.InitialDelayMs) * time.Millisecond
		for attempt := 1; attempt <= s.config.Retries.MaxAttempts; attempt++ {
			chunks, errs := s.provider.GenerateStream(ctx, input)
			hasOutput := false

			for chunks != nil || errs != nil {
				select {
				case chunk, ok := <-chunks:
					if !ok {
						chunks = nil
						continue
					}
					hasOutput = true
					select {
					case outChunks <- chunk:
					case <-ctx.Done():
						return
					}

				case streamErr, ok := <-errs:
					if !ok {
						errs = nil
						continue
					}
					if !hasOutput && !errors.Is(streamErr, context.Canceled) && !errors.Is(streamErr, context.DeadlineExceeded) && streamErrorRetryable(streamErr) && attempt < s.config.Retries.MaxAttempts {
						if !waitForStreamRetry(ctx, delay) {
							return
						}
						delay = nextRetryDelay(delay, time.Duration(s.config.Retries.MaximumDelayMs)*time.Millisecond)
						goto nextAttempt
					}
					select {
					case outErrs <- streamErr:
					case <-ctx.Done():
					}
					return
				}
			}
			return

		nextAttempt:
		}
	}()

	return outChunks, outErrs, nil
}

func streamErrorRetryable(err error) bool {
	var providerErr *provider.Error
	return errors.As(err, &providerErr) && utils.IsRetryableError(providerErr.StatusCode)
}

func waitForStreamRetry(ctx context.Context, delay time.Duration) bool {
	wait := time.Duration(rand.Int64N(int64(delay) + 1))
	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func nextRetryDelay(delay, maximum time.Duration) time.Duration {
	if delay >= maximum || delay > maximum/2 {
		return maximum
	}
	return delay * 2
}
