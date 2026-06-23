package rag

import (
	"context"
	"fmt"

	"github.com/LucasM4r/repomind/internal/cache"
	"github.com/LucasM4r/repomind/internal/db"
	"github.com/LucasM4r/repomind/internal/domain"
	"github.com/LucasM4r/repomind/internal/rpc"
)

type RAG struct {
	AIClient           *rpc.AIClient
	DB                 *db.VectorRepository
	Cache              *cache.SessionCache
	promptHistoryLimit int
}

func NewRAGService(ai *rpc.AIClient, db *db.VectorRepository, cache *cache.SessionCache, promptHistoryLimit int) *RAG {
	if limit := promptHistoryLimit; limit <= 0 {
		promptHistoryLimit = 10
	}
	return &RAG{AIClient: ai, DB: db, Cache: cache, promptHistoryLimit: promptHistoryLimit}
}

func (r *RAG) initRAGCache(ctx context.Context, sessionID string) error {
	history, err := r.DB.GetRecentHistory(ctx, sessionID, r.promptHistoryLimit)
	if err != nil {
		return fmt.Errorf("[RAG][DB] failed to get recent history: %w", err)
	}
	r.Cache.Set(sessionID, history)
	return nil
}

func (r *RAG) promptWithHistory(ctx context.Context, sessionID string, userPrompt string) (string, error) {
	history, ok := r.Cache.Get(sessionID)
	if !ok {
		fmt.Printf("[WARNING][RAG] no history found for sessionID: %s. Initializing cache...\n", sessionID)
		if err := r.initRAGCache(ctx, sessionID); err != nil {
			return "", err
		}
		history, _ = r.Cache.Get(sessionID)
	}

	userPromptMessage := domain.ChatMessage{
		Role:    "user",
		Content: userPrompt,
	}

	if err := r.DB.SaveMessage(ctx, sessionID, userPromptMessage.Role, userPromptMessage.Content); err != nil {
		return "", fmt.Errorf("[ERROR][RAG][DB] failed to save user message: %w", err)
	}

	if len(history) >= r.promptHistoryLimit {
		history = append([]domain.ChatMessage(nil), history[len(history)-r.promptHistoryLimit+1:]...)
	}
	history = append(history, userPromptMessage)
	r.updateCache(sessionID, history)

	response, err := r.AIClient.GenerateText(ctx, history)
	if err != nil {
		return "", fmt.Errorf("[ERROR][RAG][AIClient] failed to generate text: %w", err)
	}

	if err := r.DB.SaveMessage(ctx, sessionID, response.Role, response.Content); err != nil {
		return "", fmt.Errorf("[ERROR][RAG][DB] failed to save assistant message: %w", err)
	}

	history = append(history, response)
	r.updateCache(sessionID, history)
	return response.Content, nil
}

func (r *RAG) updateCache(sessionID string, history []domain.ChatMessage) {
	if len(history) > r.promptHistoryLimit {
		history = append([]domain.ChatMessage(nil), history[len(history)-r.promptHistoryLimit:]...)
	}
	r.Cache.Set(sessionID, history)
}
