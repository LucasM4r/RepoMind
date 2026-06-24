package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/LucasM4r/repomind/internal/domain"
)

func (r *RAG) retrieveRelevantChunks(
	ctx context.Context,
	owner string,
	repo string,
	userPrompt string,
	maxDistance float32,
	limit int,
) ([]domain.RetrievedChunk, error) {
	embResp, err := r.AIClient.GetEmbeddings(ctx, []string{userPrompt})
	if err != nil {
		return nil, err
	}

	if len(embResp.GetEmbeddings()) == 0 || len(embResp.GetEmbeddings()[0].GetValues()) == 0 {
		return nil, fmt.Errorf("[RAG][Retrieval] empty embedding response")
	}

	queryEmbedding := embResp.GetEmbeddings()[0].GetValues()

	chunks, err := r.DB.GetChunks(ctx, owner, repo, queryEmbedding, maxDistance, limit)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func buildRetrievedContext(chunks []domain.RetrievedChunk) string {
	if len(chunks) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("Use the following repository context when relevant:\n")

	for i, chunk := range chunks {
		builder.WriteString(fmt.Sprintf("[%d] owner=%s repo=%s distance=%.4f\n%s\n\n", i+1, chunk.Owner, chunk.Repo, chunk.Distance, chunk.Content))
	}

	return builder.String()
}
