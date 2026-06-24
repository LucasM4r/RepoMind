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

func (r *RAG) GenerateResponse(ctx context.Context, sessionID, owner, repo, userPrompt string) (string, error) {
	const retrievalMaxDistance float32 = 0.8
	const retrievalTopK = 5

	history, err := r.getHistory(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("[ERROR][RAG] failed to get history: %w", err)
	}

	chunks, err := r.retrieveRelevantChunks(ctx, owner, repo, userPrompt, retrievalMaxDistance, retrievalTopK)
	if err != nil {
		return "", fmt.Errorf("[ERROR][RAG][Retrieval] failed to retrieve relevant chunks: %w", err)
	}

	retrievedContext := buildRetrievedContext(chunks)

	userPromptMessage := domain.ChatMessage{
		Role:    "user",
		Content: userPrompt,
	}

	if err := r.DB.SaveMessage(ctx, sessionID, userPromptMessage.Role, userPromptMessage.Content); err != nil {
		return "", fmt.Errorf("[ERROR][RAG][DB] failed to save user message: %w", err)
	}

	llmHistory := trimHistory(history, r.promptHistoryLimit-1)
	if retrievedContext != "" {
		llmHistory = append(trimHistory(llmHistory, r.promptHistoryLimit-1), domain.ChatMessage{
			Role:    "system",
			Content: retrievedContext,
		})
	}
	llmHistory = append(trimHistory(llmHistory, r.promptHistoryLimit-1), userPromptMessage)

	history = append(trimHistory(history, r.promptHistoryLimit-1), userPromptMessage)
	r.updateCache(sessionID, history)

	response, err := r.AIClient.GenerateText(ctx, llmHistory)
	if err != nil {
		return "", fmt.Errorf("[ERROR][RAG][AIClient] failed to generate text: %w", err)
	}

	if err := r.DB.SaveMessage(ctx, sessionID, response.Role, response.Content); err != nil {
		return "", fmt.Errorf("[ERROR][RAG][DB] failed to save assistant message: %w", err)
	}

	history = append(trimHistory(history, r.promptHistoryLimit-1), response)
	r.updateCache(sessionID, history)
	return response.Content, nil
}
