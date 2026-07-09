package app

import (
	"context"
	"log"

	"github.com/LucasM4r/repomind/internal/cache"
	"github.com/LucasM4r/repomind/internal/chunker"
	"github.com/LucasM4r/repomind/internal/config"
	"github.com/LucasM4r/repomind/internal/db"
	"github.com/LucasM4r/repomind/internal/ingestor"
	"github.com/LucasM4r/repomind/internal/providers"
	"github.com/LucasM4r/repomind/internal/providers/github"
	"github.com/LucasM4r/repomind/internal/rag"
	"github.com/LucasM4r/repomind/internal/rpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg        *config.Config
	dbPool     *pgxpool.Pool
	grpcConn   *grpc.ClientConn
	aiClient   *rpc.AIClient
	vectorRepo *db.VectorRepository
	ingestor   *ingestor.Ingestor

	providers   map[string]providers.Fetcher
	jobsChannel chan ingestor.Job
}

func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	conn, err := grpc.NewClient(cfg.AIEngineURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	aiClient := rpc.NewAIClient(conn)

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		conn.Close()
		return nil, err
	}
	repo := db.VectorRepository{DB: pool}

	ghProvider := github.NewClient(cfg.GithubToken)
	jobsChannel := make(chan ingestor.Job, 100)

	registry := map[string]providers.Fetcher{
		"github": ghProvider,
	}

	maxChunkSize := 1500
	overlapLines := 50

	fallbackChunker := chunker.NewLineChunker(maxChunkSize, overlapLines)
	codeChunker := chunker.NewTreeSitterChunker(maxChunkSize, fallbackChunker)

	return &App{
		cfg:         cfg,
		dbPool:      pool,
		grpcConn:    conn,
		aiClient:    aiClient,
		vectorRepo:  &repo,
		providers:   registry,
		jobsChannel: jobsChannel,
		ingestor:    ingestor.NewIngestor(aiClient, &repo, codeChunker),
	}, nil
}

func (a *App) EnqueueJob(job ingestor.Job) {
	a.jobsChannel <- job
}

func (a *App) Ingestor() *ingestor.Ingestor {
	return a.ingestor
}

func (a *App) GetProvider(name string) (providers.Fetcher, bool) {
	provider, exists := a.providers[name]
	return provider, exists
}

func (a *App) Close() {
	if a.dbPool != nil {
		a.dbPool.Close()
	}

	if a.grpcConn != nil {
		a.grpcConn.Close()
	}
}

func (a *App) RAGService() *rag.RAG {
	return rag.NewRAGService(a.aiClient, a.vectorRepo, &cache.SessionCache{}, 10)
}

func (a *App) Run(ctx context.Context) {
	log.Printf("[INFO] Initializing %d background workers", a.cfg.MaxWorkers)
	for w := 1; w <= a.cfg.MaxWorkers; w++ {
		go func(workerID int) {
			log.Printf("[INFO] Worker %d started", workerID)
			ingestor.Worker(ctx, workerID, a.jobsChannel, a.ingestor)
		}(w)
	}
}
