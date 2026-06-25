package domain

import (
	"context"
)

type ContextRetriever interface {
	GetChunks(ctx context.Context, owner, repo string, embedding []float32, maxDistance float32, limit int) ([]RetrievedChunk, error)
}
