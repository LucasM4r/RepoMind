package domain

import (
	"context"
)

type ContextRetriever interface {
	GetChunks(ctx context.Context, owner, repo string, embedding []float32, maxDistance float32, limit int) ([]RetrievedChunk, error)
}

type ChatHistoryRepository interface {
	SaveMessage(ctx context.Context, sessionID, role, content string) error
	GetRecentHistory(ctx context.Context, sessionID string, limit int) ([]ChatMessage, error)
}
