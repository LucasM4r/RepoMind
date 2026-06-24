package rag

import (
	"context"
	"fmt"
	"log"

	"github.com/LucasM4r/repomind/internal/domain"
)

func (r *RAG) initRAGCache(ctx context.Context, sessionID string) error {
	history, err := r.DB.GetRecentHistory(ctx, sessionID, r.promptHistoryLimit)
	if err != nil {
		log.Printf("[RAG][DB] failed to get recent history: %v", err)
		return err
	}
	r.Cache.Set(sessionID, history)
	return nil
}

func (r *RAG) updateCache(sessionID string, history []domain.ChatMessage) {
	history = trimHistory(history, r.promptHistoryLimit)
	r.Cache.Set(sessionID, history)
}

func (r *RAG) getHistory(ctx context.Context, sessionID string) ([]domain.ChatMessage, error) {
	history, ok := r.Cache.Get(sessionID)
	if !ok {
		log.Printf("[WARNING][RAG] no history found for sessionID: %s. Initializing cache...\n", sessionID)
		if err := r.initRAGCache(ctx, sessionID); err != nil {
			return nil, fmt.Errorf("[ERROR][RAG] failed to initialize cache: %w", err)
		}
		history, _ = r.Cache.Get(sessionID)
	}
	return history, nil
}

func trimHistory(history []domain.ChatMessage, max int) []domain.ChatMessage {
	if max <= 0 {
		return nil
	}

	if len(history) > max {
		return history[len(history)-max:]
	}
	return history
}
