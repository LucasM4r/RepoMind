package domain

import "context"

type ChatHistoryRepository interface {
	SaveMessage(ctx context.Context, sessionID, role, content string) error
	GetRecentHistory(ctx context.Context, sessionID string, limit int) ([]ChatMessage, error)
}
