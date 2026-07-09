package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/LucasM4r/repomind/internal/api/handlers"
	"github.com/LucasM4r/repomind/internal/api/middleware"
	"github.com/LucasM4r/repomind/internal/rag"
)

type Router struct {
	mux     *http.ServeMux
	handler *handlers.Handler
}

func NewRouter(enqueuer handlers.JobEnqueuer, ragService *rag.RAG) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		handler: handlers.NewHandler(enqueuer, ragService),
	}

	return r
}

func (r *Router) RegisterHandlers() {
	r.mux.HandleFunc("/api/v1/ask", r.handler.AskHandler)
	r.mux.HandleFunc("/api/v1/ingest", r.handler.IngestHandler)

	r.mux.Handle("/api/docs/", httpSwagger.WrapHandler)
	r.mux.Handle("/api/docs", http.RedirectHandler("/api/docs/", http.StatusMovedPermanently))
}

func (r *Router) Handler() http.Handler {
	return middleware.CORS(r.mux)
}
