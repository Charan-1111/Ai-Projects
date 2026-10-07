package services

import (
	"context"
	"fmt"
	"time"
	"tool-calling/internal/models"
)

func (s *Service) GetCurrentTime(ctx context.Context, args models.TimeInput) (models.TimeOutput, error) {
	if args.Timezone == "" {
		return models.TimeOutput{}, fmt.Errorf("timezone is required")
	}

	location, err := time.LoadLocation(args.Timezone)
	if err != nil {
		return models.TimeOutput{}, fmt.Errorf("invalid timezone %q: %w", args.Timezone, err)
	}

	return models.TimeOutput{
		Timezone: args.Timezone,
		Datetime: time.Now().In(location).Format(time.RFC3339),
	}, nil
}
