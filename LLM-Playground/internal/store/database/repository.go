package database

import (
	"context"
	"fmt"
	"llm-playground/internal/provider"
)

type Repository interface {
	Create(ctx context.Context) error
	Messages(ctx context.Context, conversationID string) ([]provider.Message, error)
	SaveMessages(ctx context.Context, requestID, conversationID string, messages []provider.Message) error
}

func (db *DataBaseStore) Create(ctx context.Context) error {
	for tableName, createQuery := range db.Queries.Create {
		_, err := db.Db.Exec(ctx, createQuery)
		if err != nil {
			return err
		} else {
			db.Log.Log.Info().Msgf("Table %s created successfully", tableName)
		}
	}

	return nil
}

func (db *DataBaseStore) Messages(ctx context.Context, conversationID string) ([]provider.Message, error) {
	rows, err := db.Db.Query(ctx, db.Queries.Fetch.Messages, conversationID)
	if err != nil {
		return nil, fmt.Errorf("fetch conversation messages: %w", err)
	}
	defer rows.Close()

	messages := make([]provider.Message, 0)
	for rows.Next() {
		var message provider.Message
		if err := rows.Scan(&message.Role, &message.Content); err != nil {
			return nil, fmt.Errorf("scan conversation message: %w", err)
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversation messages: %w", err)
	}

	return messages, nil
}

func (db *DataBaseStore) SaveMessages(ctx context.Context, requestID, conversationID string, messages []provider.Message) error {
	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin save conversation: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		db.Queries.Save.Conversation,
		conversationID,
	)
	if err != nil {
		return fmt.Errorf("ensure conversation: %w", err)
	}

	for index, message := range messages {
		messageID := requestID
		if index > 0 {
			messageID = fmt.Sprintf("%s-%d", requestID, index)
		}
		_, err = tx.Exec(ctx,
			db.Queries.Save.Message,
			messageID, conversationID, message.Role, message.Content,
		)
		if err != nil {
			return fmt.Errorf("save conversation message: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit conversation messages: %w", err)
	}

	return nil
}
