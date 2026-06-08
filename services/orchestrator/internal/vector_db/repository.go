package vector_db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type VectorRepository struct {
	DB *pgxpool.Pool
}

func (r *VectorRepository) SaveChunk(ctx context.Context, chunkID string, content string, embedding []float32) error {
	vec := pgvector.NewVector(embedding)

	query := `INSERT INTO chunks (id, content, embedding) 
              VALUES ($1, $2, $3) 
              ON CONFLICT (id) DO UPDATE SET embedding = $3, content = $2`

	_, err := r.DB.Exec(ctx, query, chunkID, content, vec)
	return err
}
