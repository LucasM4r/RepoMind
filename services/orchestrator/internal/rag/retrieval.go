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

	const maxChunkChars = 700

	var builder strings.Builder
	builder.WriteString("Repository snippets:\n")

	for i, chunk := range chunks {
		content := chunk.Content
		if len(content) > maxChunkChars {
			content = content[:maxChunkChars] + "..."
		}

		builder.WriteString(fmt.Sprintf("Snippet %d:\n%s\n\n", i+1, content))
	}

	return builder.String()
}
