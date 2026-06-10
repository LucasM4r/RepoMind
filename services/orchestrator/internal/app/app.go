package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/LucasM4r/repomind/internal/config"
	"github.com/LucasM4r/repomind/internal/ingestor"
	"github.com/LucasM4r/repomind/internal/providers"
	"github.com/LucasM4r/repomind/internal/providers/github"
	"github.com/LucasM4r/repomind/internal/rpc"
	"github.com/LucasM4r/repomind/internal/vector_db"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg        *config.Config
	dbPool     *pgxpool.Pool
	grpcConn   *grpc.ClientConn
	aiClient   *rpc.AIClient
	vectorRepo *vector_db.VectorRepository
	ghProvider providers.Fetcher
}

func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	conn, err := grpc.NewClient(cfg.AIEngineURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	aiClient := rpc.NewAIClient(conn)

	pool, err := vector_db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		conn.Close()
		return nil, err
	}
	repo := vector_db.VectorRepository{DB: pool}

	ghProvider := github.NewClient(cfg.GithubToken)

	return &App{
		cfg:        cfg,
		dbPool:     pool,
		grpcConn:   conn,
		aiClient:   aiClient,
		vectorRepo: &repo,
		ghProvider: ghProvider,
	}, nil
}

func (a *App) Close() {
	if a.dbPool != nil {
		a.dbPool.Close()
	}

	if a.grpcConn != nil {
		a.grpcConn.Close()
	}
}

func (a *App) Run() {
	jobsChannel := make(chan ingestor.Job, 100)
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ing := ingestor.NewIngestor(a.aiClient, a.vectorRepo)
	log.Printf("[INFO] Initializing %d workers", a.cfg.MaxWorkers)
	for w := 1; w <= a.cfg.MaxWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			ingestor.Worker(ctx, workerID, jobsChannel, ing)
		}(w)
	}

	projects := []struct{ Owner, Repo string }{
		{"LucasM4r", "RT-TELEMETRY-DASHBOARD"},
	}

	for _, p := range projects {
		jobsChannel <- ingestor.Job{
			Owner:    p.Owner,
			Repo:     p.Repo,
			Provider: a.ghProvider,
		}
	}

	close(jobsChannel)
	wg.Wait()

	log.Printf("[INFO] Ingestion process finished successfully")
}
