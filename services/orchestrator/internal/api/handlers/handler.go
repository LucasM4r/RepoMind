package handlers

import "github.com/LucasM4r/repomind/internal/rag"

type Handler struct {
	rag *rag.RAG
}

func NewHandler(ragService *rag.RAG) *Handler {
	return &Handler{rag: ragService}
}
