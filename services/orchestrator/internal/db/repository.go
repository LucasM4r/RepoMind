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
	query := `INSERT INTO chunks (id, owner, repo, content, embedding) 
              VALUES ($1, $2, $3, $4, $5) 
              ON CONFLICT (id) DO UPDATE SET embedding = $5, content = $4`
	_, err := r.DB.Exec(ctx, query, chunkID, owner, repo, content, vec)
	return err
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
