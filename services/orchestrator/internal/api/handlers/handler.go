package handlers

import (
	"github.com/LucasM4r/repomind/internal/ingestor"
	"github.com/LucasM4r/repomind/internal/providers"
	"github.com/LucasM4r/repomind/internal/rag"
)

type JobEnqueuer interface {
	EnqueueJob(job ingestor.Job)
	GetProvider(name string) (providers.Fetcher, bool)
}

type Handler struct {
	enqueuer   JobEnqueuer
	ragService *rag.RAG
}

func NewHandler(enqueuer JobEnqueuer, ragService *rag.RAG) *Handler {
	return &Handler{
		enqueuer:   enqueuer,
		ragService: ragService,
	}
}
