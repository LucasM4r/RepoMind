package router

import (
	"net/http"

	"github.com/LucasM4r/repomind/internal/api/handlers"
	"github.com/LucasM4r/repomind/internal/api/middleware"
	"github.com/LucasM4r/repomind/internal/rag"
)

type Router struct {
	mux     *http.ServeMux
	handler *handlers.Handler
}

func NewRouter(ragService *rag.RAG) *Router {
	return &Router{
		mux:     http.NewServeMux(),
		handler: handlers.NewHandler(ragService),
	}
}

func (r *Router) RegisterHandlers() {
	r.mux.HandleFunc("/api/v1/ask", r.handler.AskHandler)
}

func (r *Router) Handler() http.Handler {
	return middleware.CORS(r.mux)
}
