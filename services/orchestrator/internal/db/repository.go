package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
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
