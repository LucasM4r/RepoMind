package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"

	"github.com/LucasM4r/repomind/internal/domain"
)

type VectorRepository struct {
	DB *pgxpool.Pool
}

func (r *VectorRepository) SaveChunk(ctx context.Context, owner, repo, filename, chunkID, content string, embedding []float32) error {
	vec := pgvector.NewVector(embedding)
	query := `INSERT INTO chunks (id, owner, repo, filename, content, embedding) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              ON CONFLICT (id) DO UPDATE SET embedding = $6, content = $4`
	_, err := r.DB.Exec(ctx, query, chunkID, owner, repo, filename, content, vec)
	return err
}

func (r *VectorRepository) GetChunks(
	ctx context.Context,
	owner, repo string,
	embedding []float32,
	maxDistance float32,
	limit int,
) ([]domain.RetrievedChunk, error) {
	if limit <= 0 {
		limit = 10
	}

	vec := pgvector.NewVector(embedding)
	query := `
        SELECT id, owner, repo, content, embedding <-> $3 AS distance
        FROM chunks
        WHERE owner = $1
          AND repo = $2
          AND embedding <-> $3 < $4
        ORDER BY embedding <-> $3 ASC
        LIMIT $5
    `

	rows, err := r.DB.Query(ctx, query, owner, repo, vec, maxDistance, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []domain.RetrievedChunk
	for rows.Next() {
		var chunk domain.RetrievedChunk
		if err := rows.Scan(
			&chunk.ID,
			&chunk.Owner,
			&chunk.Repo,
			&chunk.Content,
			&chunk.Distance,
		); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chunks, nil
}

func (r *VectorRepository) SaveMessage(ctx context.Context, sessionID, role, content string) error {
	query := `
        INSERT INTO chat_history (session_id, role, content) 
        VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(ctx, query, sessionID, role, content)
	return err
}

func (r *VectorRepository) GetRecentHistory(ctx context.Context, sessionID string, limit int) ([]domain.ChatMessage, error) {
	query := `SELECT id, session_id, role, content, created_at 
			FROM chat_history
			WHERE session_id = $1
			ORDER BY created_at DESC
			LIMIT $2`

	rows, err := r.DB.Query(ctx, query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.ChatMessage

	for rows.Next() {
		var msg domain.ChatMessage

		err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, msg)
	}

	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}
	return history, nil
}
